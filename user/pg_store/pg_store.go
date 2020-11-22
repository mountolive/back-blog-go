// Defines the Storage details associated with an User for Postgres
package pg_store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mountolive/back-blog-go/user/usecase"
)

var (
	ConnectionError    = errors.New("An error occurred when connecting to the DB")
	TableCreationError = errors.New("An error occurred when trying to create the table")
)

// Store implementation that uses Postgres for persistence
type PgStore struct {
	db *pgxpool.Pool
}

func NewUserPgStore(ctx context.Context, url string) (*PgStore, error) {
	db, err := pgxpool.Connect(ctx, url)
	if err != nil {
		return nil, wrapErrorInfo(ConnectionError, err.Error(), "user")
	}
	store := &PgStore{db}
	err = store.createUserTable(ctx)
	if err != nil {
		return nil, wrapErrorInfo(TableCreationError, err.Error(), "user")
	}
	return store, nil
}

func (p *PgStore) createUserTable(ctx context.Context) error {
	err := checkContext(&ctx)
	if err != nil {
		return err
	}
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	// Function for automatic setting of timestamps
	_, err = tx.Exec(ctx, `
    CREATE EXTENSION citext;

    CREATE OR REPLACE FUNCTION trigger_set_timestamp()
		RETURNS TRIGGER AS $$
		BEGIN
		  NEW.updated_at = NOW();
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;
    
	  CREATE TABLE IF NOT EXISTS users (
		  id         UUID NOT NULL PRIMARY KEY,
			email      CITEXT NOT NULL UNIQUE,
			username   VARCHAR(100) UNIQUE,
			password   VARCHAR(200) NOT NULL,
			first_name VARCHAR(100) NOT NULL,
			last_name  VARCHAR(100) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		);

		CREATE INDEX idx_email ON users (email);
    CREATE INDEX idx_username ON users (username);

		-- trigger automatic setting of timestamps
		CREATE TRIGGER set_timestamp
		BEFORE UPDATE ON users
		FOR EACH ROW
		EXECUTE PROCEDURE trigger_set_timestamp();
	`,
	)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// Creates an User and returns it (UserDto)
func (p *PgStore) Create(ctx context.Context,
	data *usecase.CreateUserDto) (*usecase.UserDto, error) {
	// TODO Implement
	return nil, nil
}

// Updates the data associated to an User
// and returns the corresponding UserDto
func (p *PgStore) Update(ctx context.Context,
	data *usecase.UpdateUserDto) (*usecase.UserDto, error) {
	// TODO Implement
	return nil, nil
}

// Updates a given User's Password
func (p *PgStore) UpdatePassword(ctx context.Context,
	data *usecase.ChangePasswordDto) error {
	// TODO Implement
	return nil
}

// Retrieves a single User from DB,
// through its Username or Email
func (p *PgStore) ReadOne(ctx context.Context,
	query *usecase.ByUsernameOrEmail) (*usecase.UserDto, error) {
	// TODO Implement
	return nil, nil
}

// Checks if User's credentials are OK
func (p *PgStore) CheckIfCorrectPassword(ctx context.Context,
	data *usecase.CheckUserAndPasswordDto) error {
	// TODO Implement
	return nil
}

func checkContext(ctx *context.Context) error {
	// TODO Implement
	return nil
}

func wrapErrorInfo(err error, msg string, store string) error {
	return fmt.Errorf("%w - %s store: %s \n", err, store, msg)
}
