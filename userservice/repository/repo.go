package repository

import (
	"context"
	"database/sql"

	"github.com/jinzhu/copier"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.uber.org/zap"

	"EchoWave/pkg/logger"
	"EchoWave/userservice/repository/postgres/models"
	"EchoWave/userservice/types"
)

type postgresRepo struct {
	conn *sql.DB
}

type Repository interface {
	// GetUserByUsername  returns a user by username
	GetUserByUsername(ctx context.Context, username string) (*types.User, error)
}

func InitPostgresRepo(db *sql.DB) Repository {
	return &postgresRepo{
		conn: db,
	}
}

func (p postgresRepo) GetUserByUsername(ctx context.Context, username string) (*types.User, error) {
	ctx, span := types.UserServiceTracer.Start(ctx, "PGRepo:GetUserByUsername")
	defer span.End()

	user, err := models.Users(qm.Where("username = ?", username)).One(ctx, p.conn)
	if err != nil {
		logger.Instance().Ctx(ctx).Info("error getting user by username", zap.Error(err), zap.String("username", username))

		return nil, err
	}

	resp := new(types.User)
	err = copier.Copy(resp, user)

	return resp, err
}
