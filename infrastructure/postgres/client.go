package postgres

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"github.com/volatiletech/sqlboiler/v4/boil"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.uber.org/zap"

	"EchoWave/pkg/logger"
)

func Dial(url string) (*sql.DB, error) {
	db, err := otelsql.Open("pgx", url,
		otelsql.WithAttributes(semconv.DBSystemPostgreSQL),
		otelsql.WithDBName("EchoWave/postgres"),
	)
	if err != nil {
		logger.Instance().Error("error opening database", zap.Error(err))

		return nil, err
	}

	err = db.Ping()
	if err != nil {
		logger.Instance().Error("error pinging database", zap.Error(err))

		return nil, err
	}

	boil.SetDB(db)

	return db, nil
}

func Cleanup(db *sql.DB) {
	if err := db.Close(); err != nil {
		logger.Instance().Error("error closing database", zap.Error(err))
	}
}
