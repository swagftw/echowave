package postgres

import (
	"database/sql"

	"go.uber.org/zap"

	"EchoWave/pkg/logger"
)

func Migrate(conn *sql.DB) error {
	// TODO: use migration library
	// now using simple sql statements with pgx

	// create user schema if it does not exist
	query := "CREATE SCHEMA IF NOT EXISTS usr"

	_, err := conn.Exec(query)
	if err != nil {
		logger.Instance().Error("error creating user schema", zap.Error(err))

		return err
	}

	// create user table if it does not exist
	query = `CREATE TABLE IF NOT EXISTS usr.users (
    			id SERIAL PRIMARY KEY,
    			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
                deleted_at TIMESTAMP,
                username VARCHAR(255) UNIQUE NOT NULL,
                password VARCHAR(255) NOT NULL
			)`

	_, err = conn.Exec(query)
	if err != nil {
		logger.Instance().Error("error creating user table", zap.Error(err))

		return err
	}

	return nil
}
