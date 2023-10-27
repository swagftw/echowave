package server

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"EchoWave/gateway/config"
	"EchoWave/gateway/routes/user"
	"EchoWave/internal/proto"
	"EchoWave/pkg/logger"
)

// CreateServer creates a new fiber API gateway server.
func CreateServer(cfg *config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.IdleTimeout) * time.Second,
	})

	// middlewares
	// using default configurations for now
	app.Use(recover.New(), fiberLogger.New(), requestid.New())
	app.Use(otelfiber.Middleware())

	return app
}

// InitGateway initializes the gateway.
func InitGateway(app *fiber.App, cfg *config.Config) {
	ctx := context.Background()

	// initialize uptrace for tracing
	uptrace.ConfigureOpentelemetry(
		uptrace.WithDSN(cfg.OTel.DSN),
		uptrace.WithServiceName(cfg.OTel.Service),
		uptrace.WithServiceVersion(cfg.OTel.Version),
		uptrace.WithDeploymentEnvironment(cfg.OTel.DeploymentEnv),
	)

	defer func(ctx context.Context) {
		err := uptrace.Shutdown(ctx)
		if err != nil {
			slog.Error("error shutting down uptrace", "error", err)
		}
	}(ctx)

	logger.Instance().Info("initializing gateway", zap.Int("port", cfg.Port))

	v1Route := app.Group("/api/v1")

	// initial user service grpc client
	userService, err := initUserGRPCClient(cfg)
	if err != nil {
		return
	}

	// initialize user routes
	user.InitRoutes(v1Route, userService)

	hostPort := net.JoinHostPort("0.0.0.0", strconv.Itoa(cfg.Port))

	// errChan is for handling errors from the server
	errChan := make(chan error)
	// sigChan is for handling signals from the OS
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		err := app.Listen(hostPort)
		if err != nil {
			errChan <- err
		}
	}()

	// block until an error or signal is received
	select {
	case err := <-errChan:
		logger.Instance().Error("error starting server", zap.Error(err))

		return
	case sig := <-sigChan:
		logger.Instance().Info("received signal", zap.String("signal", sig.String()))

		return
	}
}

func initUserGRPCClient(cfg *config.Config) (*proto.UserServiceClient, error) {
	conn, err := grpc.Dial(cfg.UserService.GrpcURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		logger.Instance().Error("error dialing user service grpc server", zap.Error(err))

		return nil, err
	}

	client := proto.NewUserServiceClient(conn)

	return &client, nil
}
