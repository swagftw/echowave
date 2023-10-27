package postgres

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.uber.org/zap"

	"EchoWave/pkg/logger"
)

func Dial(url string) (*sql.DB, error) {
	db, err := sql.Open("pgx", url)
	if err != nil {
		logger.Instance().Error("error opening database", zap.Error(err))

		return nil, err
	}

	boil.SetDB(db)

	return db, err
}

func Cleanup(db *sql.DB) {
	if err := db.Close(); err != nil {
		logger.Instance().Error("error closing database", zap.Error(err))
	}
}
