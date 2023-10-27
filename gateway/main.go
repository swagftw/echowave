package main

import (
	"flag"
	"log/slog"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	gatewayConfig "EchoWave/gateway/config"
	"EchoWave/gateway/server"
	"EchoWave/pkg/logger"
)

func main() {
	defer func(instance *otelzap.Logger) {
		err := instance.Sync()
		if err != nil {
			slog.Error("error syncing logger", "error", err)
		}
	}(logger.Instance())

	// read config for api gateway
	gatewayConfigPath := flag.String("gateway-config", "./gateway/config.yaml", "path to config file")
	flag.Parse()

	cfg, err := gatewayConfig.ReadConfig(*gatewayConfigPath)
	if err != nil {
		return
	}

	// create server
	svr := server.CreateServer(cfg)

	// run server
	server.InitGateway(svr, cfg)
}
