package main

import (
	"context"
	"flag"

	"github.com/uptrace/uptrace-go/uptrace"
	"go.uber.org/zap"

	"EchoWave/infrastructure/postgres"
	"EchoWave/pkg/logger"
	"EchoWave/userservice/config"
	"EchoWave/userservice/repository"
	userPg "EchoWave/userservice/repository/postgres"
	"EchoWave/userservice/server"
	"EchoWave/userservice/service"
)

func main() {
	configPath := flag.String("config", "./userservice/config.yaml", "path to config file")

	ctx := context.Background()

	cfg, err := config.ReadConfig(*configPath)
	if err != nil {
		return
	}

	// initialize uptrace for tracing
	uptrace.ConfigureOpentelemetry(
		uptrace.WithDSN(cfg.OTel.DSN),
		uptrace.WithServiceName(cfg.OTel.Service),
		uptrace.WithServiceVersion(cfg.OTel.Version),
		uptrace.WithDeploymentEnvironment(cfg.OTel.DeploymentEnv),
	)

	conn, err := postgres.Dial(cfg.Postgres.URL)
	if err != nil {
		return
	}

	defer func(ctx context.Context) {
		err = uptrace.Shutdown(ctx)
		if err != nil {
			logger.Instance().Error("error shutting down uptrace", zap.Error(err))
		}

		postgres.Cleanup(conn)
	}(ctx)

	err = userPg.Migrate(conn)
	if err != nil {
		return
	}

	// init repository
	postgresRepo := repository.InitPostgresRepo(conn)

	// init service
	userService := service.InitService(postgresRepo)

	server.StartServer(userService, cfg)
}
