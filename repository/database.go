package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var DB DBProvider

type DBProvider struct {
	Connection *sqlx.DB
}

func InitDB(typeDB, conn string) {
	var (
		err error
	)

	DB.Connection, err = sqlx.Connect(typeDB, conn)
	if err != nil {
		log.Fatalf("[InitDB][sqlx.Connect] Input: %v Output: %v", conn, err)
	}

	if _, err = DB.QueryContext(context.Background(), "SELECT 1 FROM user_credential", nil); err != nil {
		log.Fatalf("[InitDB][QueryContext] Input: %v Output: %v", conn, err)
	}
}

func (db *DBProvider) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return db.Connection.QueryRowContext(ctx, query, args...)
}

func (db *DBProvider) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return db.Connection.QueryContext(ctx, query, args...)
}

func (db *DBProvider) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return db.Connection.ExecContext(ctx, query, args...)
}
