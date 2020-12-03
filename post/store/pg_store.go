// Defines the corresponding store associated to posts, with Postgres
package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mountolive/back-blog-go/post/usecase"
)

type PgStore struct {
	db *pgxpool.Pool
}

var (
	ConnectionError    = errors.New("An error occurred when connecting to the DB")
	TableCreationError = errors.New("An error occurred when trying to create needed tables")
)

func NewPostPgStore(ctx context.Context, url string) (*PgStore, error) {
	db, err := pgxpool.Connect(ctx, url)
	if err != nil {
		return nil, wrapErrorInfo(ConnectionError, err.Error(), "post")
	}

	store := &PgStore{db}
	err = store.createPostAndTagTable(ctx)
	if err != nil {
		return nil, wrapErrorInfo(TableCreationError, err.Error(), "post")
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
        content    VARCHAR(5000) CHECK (content <> ''),
        created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
      );

      CREATE INDEX IF NOT EXISTS idx_creator ON posts (creator);

      CREATE TABLE IF NOT EXISTS tags (
        id         UUID NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
        name       CITEXT NOT NULL UNIQUE CHECK (name <> '')
      );

      CREATE INDEX IF NOT EXISTS idx_name ON tags (name);

      CREATE TABLE IF NOT EXISTS posts_tags (
        post_id    UUID NOT NULL,
        tag_id     UUID NOT NULL,
        CONSTRAINT fk_post FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
        CONSTRAINT fk_tag FOREIGN KEY (tag_id) REFERENCES tags (id) ON DELETE CASCADE
      );

      CREATE INDEX IF NOT EXISTS idx_post_id ON posts_tags (post_id);
      CREATE INDEX IF NOT EXISTS idx_tag_id ON posts_tags (tag_id);

      DROP TRIGGER IF EXISTS set_timestamp ON posts;
      CREATE TRIGGER set_timestamp
      BEFORE UPDATE ON posts
      FOR EACH ROW
      EXECUTE PROCEDURE trigger_set_timestamp();
  `)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (p *PgStore) Create(ctx context.Context,
	create *usecase.CreatePostDto) (*usecase.PostDto, error) {
	// TODO implement
	return nil, nil
}

func (p *PgStore) Update(ctx context.Context,
	update *usecase.UpdatePostDto) (*usecase.PostDto, error) {
	// TODO implement
	return nil, nil
}

func (p *PgStore) Filter(ctx context.Context,
	filter *usecase.GeneralFilter) ([]*usecase.PostDto, error) {
	// TODO implement
	return []*usecase.PostDto{}, nil
}

func wrapErrorInfo(err error, msg string, store string) error {
	return fmt.Errorf("%w - %s store: %s \n", err, store, msg)
}
