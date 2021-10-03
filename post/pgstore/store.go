// Defines the corresponding store associated to posts, with Postgres
package pgstore

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/mountolive/back-blog-go/post/usecase"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

// Store wrapper, it uses a connection pool
type PgStore struct {
	db *pgxpool.Pool
}

var (
	// TODO Add documentation and rename post's store errors
	ConnectionError        = errors.New("error occurred when connecting to the DB")
	TableCreationError     = errors.New("error occurred when trying to create needed tables")
	CreateTransactionError = errors.New("error occurred when trying to create transaction")
	ExecTransactionError   = errors.New("error occurred when trying to exec transaction")
	FilterError            = errors.New("error occurred when trying to exec Filter query")
	ReadOneError           = errors.New("error occurred when trying to exec Read query")
)

const (
	insertTag = `
         INSERT INTO tags (tag_name) VALUES %s
         ON CONFLICT (tag_name) DO UPDATE SET tag_name = EXCLUDED.tag_name RETURNING id
  `
	selectPost = `
         SELECT
           id, p.creator, p.title, p.content, p.created_at, p.updated_at, t.tag_array
         FROM posts p LEFT OUTER JOIN (
           SELECT pt.post_id AS id, array_agg(tg.tag_name)::text[] AS tag_array
           FROM posts_tags pt
           JOIN tags tg ON tg.id = pt.tag_id
           GROUP BY pt.post_id
         ) t USING (id)
         %s;
  `
	insertPostTag = `
         INSERT INTO posts_tags (post_id, tag_id)
         SELECT p.id, t.id FROM tagids t CROSS JOIN postids p
         ON CONFLICT (post_id, tag_id)
         DO UPDATE SET post_id=EXCLUDED.post_id, tag_id=EXCLUDED.tag_id;
  `
)

// Creates a store for persistences of posts
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

// Creates a Post with data with corresponding CreatePostDto
// TODO Remove Post from return tuple, Create, pgstore
func (p *PgStore) Create(ctx context.Context,
	create *usecase.CreatePostDto) (*usecase.Post, error) {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return nil, wrapErrorInfo(CreateTransactionError, err.Error())
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	postStatement := "INSERT INTO posts (creator, title, content) VALUES ($1, $2, $3) RETURNING id"
	insertTagStatement := fmt.Sprintf(insertTag, insertParamsString(create.Tags, 4))
	joinUpsert := `
         WITH postids AS (
           %s
         ),
         tagids AS (
           %s
         )

         %s
  `
	statement := fmt.Sprintf(
		joinUpsert,
		postStatement, insertTagStatement, insertPostTag,
	)
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

	// TODO Post's read by Creator in Create's return is an open door for inconsistent results (write skew)
	return p.getNewestByCreator(ctx, create.Creator), nil
}

// Updates the corresponding post with the Id from the UpdatePostDto passed
// TODO Remove Post from return tuple, Update, pgstore
func (p *PgStore) Update(ctx context.Context,
	update *usecase.UpdatePostDto) (*usecase.Post, error) {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return nil, wrapErrorInfo(CreateTransactionError, err.Error())
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()
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
	// TODO Post's ReadOne in Update's return is an open door for inconsistent results (write skew)
	return p.ReadOne(ctx, update.Id)
}

// Reads from the store the post with the passed Id
func (p *PgStore) ReadOne(ctx context.Context, id string) (*usecase.Post, error) {
	post := &usecase.Post{}
	row := p.db.QueryRow(
		ctx,
		fmt.Sprintf(selectPost, "WHERE id = $1"),
		id,
	)
	err := rowToPost(row, post)
	if err != nil {
		return nil, wrapErrorInfo(ReadOneError, err.Error())
	}
	return post, nil
}

// Filters either by tags and/or creation date
func (p *PgStore) Filter(ctx context.Context,
	filter *usecase.GeneralFilter) ([]*usecase.Post, error) {
	params := make([]interface{}, 0)
	filterStatement := buildFilterStatement(filter, &params)
	rows, err := p.db.Query(ctx, filterStatement, params...)
	if err != nil {
		return nil, wrapErrorInfo(FilterError, err.Error())
	}
	defer rows.Close()
	posts := []*usecase.Post{}
	for rows.Next() {
		post := &usecase.Post{}
		err = rowToPost(rows, post)
		if err != nil {
			return nil, wrapErrorInfo(FilterError, err.Error())
		}
		posts = append(posts, post)
	}
	if rows.Err() != nil {
		return nil, wrapErrorInfo(FilterError, err.Error())
	}
	return posts, nil
}

// CreateTestContainer creates a DB container for integration tests
func CreateTestContainer(t *testing.T, containerName string) *PgStore {
	err := godotenv.Load("../.env.test")
	if err != nil {
		t.Log("dotenv file not found")
	}
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Error starting docker: %s\n", err)
	}

	testPass := os.Getenv("POSTGRES_TEST_PASSWORD")
	testUser := os.Getenv("POSTGRES_TEST_USER")
	hostPort := "5434"
	containerPort := "5432"

	runOptions := &dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13.1",
		Env: []string{
			fmt.Sprintf("POSTGRES_PASSWORD=%s", testPass),
			fmt.Sprintf("POSTGRES_USER=%s", testUser),
		},
		ExposedPorts: []string{containerPort},
		Name:         containerName,
		PortBindings: map[docker.Port][]docker.PortBinding{
			docker.Port(containerPort): {
				{
					HostIP:   "0.0.0.0",
					HostPort: hostPort,
				},
			},
		},
	}

	container, err := pool.RunWithOptions(runOptions)
	if err != nil {
		t.Fatalf("error occurred setting the container: %s", err)
	}
	t.Cleanup(func() {
		if err := pool.Purge(container); err != nil {
			err2 := pool.RemoveContainerByName(containerName)
			if err2 != nil {
				t.Fatalf("error removing the container: %s\n", err2)
			}
			_ = pool.RemoveContainerByName(containerName)
			t.Fatalf("error purging the container: %s\n", err)
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var store *PgStore
	retryFunc := func() error {
		store, err = NewPostPgStore(ctx,
			fmt.Sprintf("postgresql://%s:%s@localhost:%s/%s?sslmode=disable",
				testUser, testPass, hostPort, testUser),
		)
		return err
	}
	if err := pool.Retry(retryFunc); err != nil {
		errRemove := pool.RemoveContainerByName(containerName)
		if errRemove != nil {
			t.Fatalf("error removing the container: %s\n", errRemove)
		}
		t.Fatalf("An error occurred initializing the db: %s\n", err)
	}
	return store
}

func (p *PgStore) createPostAndTagTable(ctx context.Context) error {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()
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
           creator    TEXT CHECK (creator <> ''),
           title      TEXT CHECK (title <> ''),
           content    TEXT CHECK (content <> ''),
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

func (p *PgStore) getNewestByCreator(ctx context.Context,
	creator string) *usecase.Post {
	post := &usecase.Post{}
	row := p.db.QueryRow(
		ctx,
		fmt.Sprintf(selectPost, "WHERE p.creator = $1 ORDER BY updated_at DESC LIMIT 1"),
		creator,
	)
	err := rowToPost(row, post)
	if err != nil {
		return nil
	}
	return post
}

func buildFilterStatement(
	filter *usecase.GeneralFilter, params *[]interface{},
) string {
	statementIdx := 1
	whereClauseSegments := []string{}
	if filter.Tag != "" {
		whereClauseSegments = append(
			whereClauseSegments,
			fmt.Sprintf(`
         id IN (
           SELECT pt.post_id FROM posts_tags pt
           LEFT OUTER JOIN tags t ON t.id = pt.tag_id
           WHERE t.tag_name	= $%d
         )
      `,
				statementIdx,
			),
		)
		*params = append(*params, filter.Tag)
		statementIdx += 1
	}
	if !filter.From.IsZero() {
		whereClauseSegments = append(
			whereClauseSegments,
			fmt.Sprintf("p.created_at >= $%d", statementIdx),
		)
		*params = append(*params, filter.From)
		statementIdx += 1
	}
	if !filter.To.IsZero() {
		whereClauseSegments = append(
			whereClauseSegments,
			fmt.Sprintf("p.created_at <= $%d", statementIdx),
		)
		*params = append(*params, filter.To)
	}
	var whereClause string
	if len(*params) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauseSegments, " AND ")
	}
	whereClause += fmt.Sprintf(
		" ORDER BY created_at DESC LIMIT %d OFFSET %d",
		filter.PageSize,
		filter.Page,
	)
	return fmt.Sprintf(selectPost, whereClause)
}

func buildUpdateStatement(update *usecase.UpdatePostDto) string {
	separated := []string{}
	preparedIndex := 1
	deleteAndJoinUpsert := `
         WITH deleted AS (
           %s
         ),
         postids AS (
           %s
         ),
         tagids AS (
           %s
         )

         %s
  `
	deleteStatement := deleteOldTagsStatement(&preparedIndex)
	checkAndAppendAssignment(
		update.Content, "content", &separated, &preparedIndex,
	)
	checkAndAppendAssignment(
		update.Title, "title", &separated, &preparedIndex,
	)
	updateStatement := fmt.Sprintf(
		"UPDATE posts SET %s WHERE id = $%d RETURNING id",
		strings.Join(separated, ","),
		preparedIndex,
	)
	insertTagStatement := fmt.Sprintf(
		insertTag,
		insertParamsString(update.Tags, preparedIndex+1),
	)
	tagPostStatement := fmt.Sprintf(
		deleteAndJoinUpsert,
		deleteStatement, updateStatement,
		insertTagStatement, insertPostTag,
	)
	return tagPostStatement
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

func checkAndAppendAssignment(param, paramName string,
	separated *[]string, index *int) {
	if param != "" {
		*separated = append(
			*separated, fmt.Sprintf("%s = $%d", paramName, *index),
		)
		*index += 1
	}
}

func deleteOldTagsStatement(index *int) string {
	statement := fmt.Sprintf(`DELETE FROM posts_tags WHERE post_id = $%d`, *index)
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

func rowToPost(rawPost pgx.Row, post *usecase.Post) error {
	return rawPost.Scan(
		&post.Id, &post.Creator,
		&post.Title, &post.Content,
		&post.CreatedAt, &post.UpdatedAt,
		&post.Tags,
	)
}

func wrapErrorInfo(err error, msg string) error {
	return fmt.Errorf("POST STORE: %w - %s\n", err, msg)
}
