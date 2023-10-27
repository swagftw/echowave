package server

import (
	"context"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"EchoWave/internal/proto"
	"EchoWave/pkg/logger"
	"EchoWave/userservice/config"
	"EchoWave/userservice/types"
)

type server struct {
	proto.UnimplementedUserServiceServer
	userService types.Service
}

func (s server) GetUserByUsername(ctx context.Context, req *proto.GetUserByUsernameRequest) (*proto.User, error) {
	user, err := s.userService.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	return &proto.User{
		Id:       user.ID,
		Username: user.Username,
		Avatar:   user.Avatar,
	}, nil
}

func StartServer(userService types.Service, cfg *config.Config) {
	hostPort := net.JoinHostPort("0.0.0.0", strconv.Itoa(cfg.Port))

	lis, err := net.Listen("tcp", hostPort)
	if err != nil {
		logger.Instance().Fatal("failed to create listener")
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)

	userServiceServer := &server{userService: userService}

	proto.RegisterUserServiceServer(grpcServer, userServiceServer)

	// errChan is for handling errors from the server
	errChan := make(chan error)
	// sigChan is for handling signals from the OS
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Instance().Info("starting user service grpc server", zap.String("hostPort", hostPort))

		err = grpcServer.Serve(lis)
		if err != nil {
			logger.Instance().Error("failed to start user service grpc server", zap.Error(err))
			errChan <- err
		}
	}()

	// block until an error or signal is received
	select {
	case err = <-errChan:
		logger.Instance().Error("failed to start user service grpc server", zap.Error(err))

		return
	case sig := <-sigChan:
		logger.Instance().Info("received signal", zap.String("signal", sig.String()))

		return
	}
}
