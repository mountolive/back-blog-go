// Defines the corresponding store associated to posts, with Postgres
package store

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mountolive/back-blog-go/post/usecase"
)

type PgStore struct {
	db *pgxpool.Pool
}

var (
	ConnectionError        = errors.New("An error occurred when connecting to the DB")
	TableCreationError     = errors.New("An error occurred when trying to create needed tables")
	CreateTransactionError = errors.New(
		"An error occurred when trying to create Create transaction",
	)
	ExecTransactionError = errors.New("An error occurred when trying to exec Create transaction")
)

const (
	insertTag = `
        INSERT INTO tags (tag_name) VALUES %s
        ON CONFLICT (tag_name) DO NOTHING RETURNING id
`

	joinUpsert = `
         WITH postids AS (
           %s
         ),
         tagids AS (
           %s
         )

         INSERT INTO posts_tags (post_id, tag_id)
         SELECT p.id, t.id FROM tagids t CROSS JOIN postids p;
`
)

func NewPostPgStore(ctx context.Context, url string) (*PgStore, error) {
	db, err := pgxpool.Connect(ctx, url)
	if err != nil {
		return nil, wrapErrorInfo(ConnectionError, err.Error())
	}

	store := &PgStore{db}
	err = store.createPostAndTagTable(ctx)
	if err != nil {
		return nil, wrapErrorInfo(TableCreationError, err.Error())
	}

	return store, nil
}

func (p *PgStore) createPostAndTagTable(ctx context.Context) error {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
        CREATE EXTENSION IF NOT EXISTS citext;
        CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

         CREATE OR REPLACE FUNCTION trigger_set_timestamp()
         RETURNS TRIGGER AS $$
         BEGIN
           NEW.updated_at = NOW();
           RETURN NEW;
         END;
         $$ LANGUAGE plpgsql;

         CREATE TABLE IF NOT EXISTS posts (
           id         UUID NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
           creator    VARCHAR(100) CHECK (creator <> ''),
           title      VARCHAR(500) CHECK (title <> ''),
           content    VARCHAR(5000) CHECK (content <> ''),
           created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
           updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
         );

         CREATE INDEX IF NOT EXISTS idx_creator ON posts (creator);

         CREATE TABLE IF NOT EXISTS tags (
           id         UUID NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
           tag_name   CITEXT NOT NULL UNIQUE CHECK (tag_name <> '')
         );

         CREATE INDEX IF NOT EXISTS idx_name ON tags (tag_name);

         CREATE TABLE IF NOT EXISTS posts_tags (
           post_id    UUID NOT NULL,
           tag_id     UUID NOT NULL,
           CONSTRAINT fk_post FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
           CONSTRAINT fk_tag FOREIGN KEY (tag_id) REFERENCES tags (id) ON DELETE CASCADE,
           CONSTRAINT post_tag_id PRIMARY KEY (post_id, tag_id)
         );

         CREATE INDEX IF NOT EXISTS idx_post_id ON posts_tags (post_id);
         CREATE INDEX IF NOT EXISTS idx_tag_id ON posts_tags (tag_id);

         DROP TRIGGER IF EXISTS set_timestamp ON posts;
         CREATE TRIGGER set_timestamp
         BEFORE UPDATE ON posts
         FOR EACH ROW
         EXECUTE PROCEDURE trigger_set_timestamp();`,
	)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (p *PgStore) Create(ctx context.Context,
	create *usecase.CreatePostDto) (*usecase.PostDto, error) {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return nil, wrapErrorInfo(CreateTransactionError, err.Error())
	}
	defer tx.Rollback(ctx)

	postStatement := "INSERT INTO posts (creator, title, content) VALUES ($1, $2, $3) RETURNING id"
	insertTagStatement := fmt.Sprintf(insertTag, insertParamsString(create.Tags, 4))

	statement := fmt.Sprintf(joinUpsert, postStatement, insertTagStatement)

	params := make([]interface{}, 0)
	params = append(params, create.Creator)
	params = append(params, create.Title)
	params = append(params, create.Content)
	for _, tag := range create.Tags {
		params = append(params, tag)
	}

	_, err = tx.Exec(ctx, statement, params...)
	if err != nil {
		return nil, wrapErrorInfo(ExecTransactionError, err.Error())
	}
	err = tx.Commit(ctx)
	if err != nil {
		return nil, wrapErrorInfo(ExecTransactionError, err.Error())
	}

	return p.getNewestByCreator(ctx, create.Creator), nil
}

func (p *PgStore) Update(ctx context.Context,
	update *usecase.UpdatePostDto) (*usecase.PostDto, error) {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return nil, wrapErrorInfo(CreateTransactionError, err.Error())
	}
	defer tx.Rollback(ctx)

	statement := buildUpdateStatement(update)
	params := buildUpdateParams(update)

	_, err = tx.Exec(ctx, statement, params...)
	if err != nil {
		return nil, wrapErrorInfo(ExecTransactionError, err.Error())
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, wrapErrorInfo(ExecTransactionError, err.Error())
	}
	return p.ReadOne(ctx, update.Id), nil
}

func (p *PgStore) ReadOne(ctx context.Context, id string) *usecase.PostDto {
	// TODO Add tags
	post := &usecase.PostDto{}
	row := p.db.QueryRow(ctx, `
         SELECT id, creator, title, content, created_at, updated_at FROM posts
         WHERE id = $1;`, id)
	rowToPost(row, post)
	return post
}

func (p *PgStore) Filter(ctx context.Context,
	filter *usecase.GeneralFilter) ([]*usecase.PostDto, error) {
	// TODO implement
	return []*usecase.PostDto{}, nil
}

func (p *PgStore) getNewestByCreator(ctx context.Context,
	creator string) *usecase.PostDto {
	// TODO Add tags
	post := &usecase.PostDto{}
	row := p.db.QueryRow(ctx, `
         SELECT id, creator, title, content, created_at, updated_at FROM posts
         WHERE creator = $1 ORDER BY updated_at DESC LIMIT 1;`, creator)
	rowToPost(row, post)
	return post
}

func buildUpdateStatement(update *usecase.UpdatePostDto) string {
	separated := []string{}
	preparedIndex := 1

	deleteStatement := deleteOldTagsStatement(&preparedIndex)

	checkAndAppendAssignment(
		update.Content, "content", &separated, &preparedIndex,
	)
	checkAndAppendAssignment(
		update.Title, "title", &separated, &preparedIndex,
	)

	updateStatement := fmt.Sprintf(
		"UPDATE posts SET %s WHERE id = $%d",
		strings.Join(separated, ","),
		preparedIndex,
	)
	fmt.Println(updateStatement)

	insertTagStatement := fmt.Sprintf(
		insertTag,
		insertParamsString(update.Tags, preparedIndex),
	)

	tagPostStatement := fmt.Sprintf(joinUpsert, updateStatement, insertTagStatement)

	return fmt.Sprintf("%s; %s", deleteStatement, tagPostStatement)
}

func buildUpdateParams(update *usecase.UpdatePostDto) []interface{} {
	params := make([]interface{}, 0)
	id := update.Id
	params = append(params, id)
	content := update.Content
	if content != "" {
		params = append(params, content)
	}
	title := update.Title
	if title != "" {
		params = append(params, title)
	}
	params = append(params, id)
	for _, tag := range update.Tags {
		params = append(params, tag)
	}
	return params
}

func checkAndAppendAssignment(
	param, paramName string, separated *[]string, index *int) {
	if param != "" {
		*separated = append(
			*separated, fmt.Sprintf("%s = $%d", paramName, *index),
		)
		*index += 1
	}
}

func deleteOldTagsStatement(index *int) string {
	statement := fmt.Sprintf(`DELETE FROM posts_tags WHERE post_id = $%d;`, *index)
	*index += 1
	return statement
}

func insertParamsString(tags []string, position int) string {
	statements := []string{}
	for i := 0; i < len(tags); i++ {
		row := fmt.Sprintf("($%d)", position)
		position++
		statements = append(statements, row)
	}
	return strings.Join(statements, ",")
}

func rowToPost(rawPost pgx.Row, post *usecase.PostDto) {
	rawPost.Scan(
		&post.Id, &post.Creator,
		&post.Title, &post.Content,
		&post.CreatedAt, &post.UpdatedAt,
	)
}

func wrapErrorInfo(err error, msg string) error {
	return fmt.Errorf("POST STORE: %w - %s\n", err, msg)
}
