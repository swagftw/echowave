package service

import (
	"context"

	"EchoWave/userservice/repository"
	"EchoWave/userservice/types"
)

type service struct {
	repo repository.Repository
}

func InitService(repo repository.Repository) types.Service {
	return &service{repo: repo}
}

func (s service) GetUserByUsername(ctx context.Context, username string) (*types.User, error) {
	ctx, span := types.UserServiceTracer.Start(ctx, "GetUserByUsername")
	defer span.End()

	return s.repo.GetUserByUsername(ctx, username)
}
