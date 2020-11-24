// Defines the Storage details associated with an User
package store

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mountolive/back-blog-go/user/usecase"
	"golang.org/x/crypto/bcrypt"
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
	err := checkContext(ctx)
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
    CREATE EXTENSION IF NOT EXISTS citext;
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

    CREATE OR REPLACE FUNCTION trigger_set_timestamp()
		RETURNS TRIGGER AS $$
		BEGIN
		  NEW.updated_at = NOW();
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;
    
	  CREATE TABLE IF NOT EXISTS users (
		  id         UUID NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
			email      CITEXT NOT NULL UNIQUE CHECK (email <> ''),
			username   VARCHAR(100) UNIQUE CHECK (username <> ''),
			password   VARCHAR(200) NOT NULL CHECK (password <> ''),
			first_name VARCHAR(100) NOT NULL CHECK (first_name <> ''),
			last_name  VARCHAR(100) NOT NULL CHECK (last_name <> ''),
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		);

		CREATE INDEX IF NOT EXISTS idx_email ON users (email);
    CREATE INDEX IF NOT EXISTS idx_username ON users (username);

		-- trigger automatic setting of timestamps
		DROP TRIGGER IF EXISTS set_timestamp ON users;
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
	err := checkContext(ctx)
	if err != nil {
		return nil, err
	}

	columns := "(email, password, first_name, last_name"
	params := "($1, $2, $3, $4"
	closingColumn := ")"
	closingParams := ")"
	if data.Username != "" {
		closingColumn = ", username)"
		closingParams = ", $5)"
	}
	statement := fmt.Sprintf(`
    INSERT INTO users %s
		VALUES %s;
  `, columns+closingColumn, params+closingParams)

	tx, err := p.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(data.Password), 10)
	if err != nil {
		return nil, err
	}
	_, err = tx.Exec(ctx, statement, data.Email, string(hashedPass),
		data.FirstName, data.LastName, data.Username)
	if err != nil {
		return nil, err
	}

	rawUser := userByEmailStatement(ctx, tx, data.Email)
	tx.Commit(ctx)
	user := &usecase.UserDto{}
	rowToEntity(rawUser, user)
	return user, nil
}

// Updates the data associated to an User
// and returns the corresponding UserDto
func (p *PgStore) Update(ctx context.Context, id string,
	data *usecase.UpdateUserDto) (*usecase.UserDto, error) {
	err := checkContext(ctx)
	if err != nil {
		return nil, err
	}
	updates := []string{}
	level := 1
	params := make([]interface{}, 0)

	if data.Email != "" {
		updates = append(updates, fmt.Sprintf("email = $%d", level))
		params = append(params, data.Email)
		level += 1
	}
	if data.Username != "" {
		updates = append(updates, fmt.Sprintf("username = $%d", level))
		params = append(params, data.Username)
		level += 1
	}
	if data.FirstName != "" {
		updates = append(updates, fmt.Sprintf("first_name = $%d", level))
		params = append(params, data.FirstName)
		level += 1
	}
	if data.LastName != "" {
		updates = append(updates, fmt.Sprintf("last_name = $%d", level))
		params = append(params, data.LastName)
		level += 1
	}
	params = append(params, id)

	tx, err := p.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, fmt.Sprintf(`
    UPDATE users
		SET %s
		WHERE id = $%d;
  `, strings.Join(updates, ", "), level), params...)
	if err != nil {
		return nil, err
	}

	rawUser := userByIdStatement(ctx, tx, id)
	user := &usecase.UserDto{}
	tx.Commit(ctx)
	rowToEntity(rawUser, user)
	return user, nil
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

func checkContext(ctx context.Context) error {
	// TODO Implement
	return nil
}

func userByEmailStatement(ctx context.Context,
	tx pgx.Tx, email string) pgx.Row {
	return tx.QueryRow(ctx, `
    SELECT
		  id, email, first_name,
			last_name, username, created_at, updated_at
		FROM users WHERE email = $1;
	`, email)
}

func userByIdStatement(ctx context.Context, tx pgx.Tx, id string) pgx.Row {
	return tx.QueryRow(ctx, `
    SELECT
		  id, email, first_name,
			last_name, username, created_at, updated_at
		FROM users WHERE id = $1;
	`, id)
}

func rowToEntity(rawUser pgx.Row, user *usecase.UserDto) {
	rawUser.Scan(&user.Id, &user.Email, &user.FirstName,
		&user.LastName, &user.Username,
		&user.CreatedAt, &user.UpdatedAt)
}

func wrapErrorInfo(err error, msg string, store string) error {
	return fmt.Errorf("%w - %s store: %s \n", err, store, msg)
}
