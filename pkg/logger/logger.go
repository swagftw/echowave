package logger

import (
	"os"
	"sync"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

var (
	once   sync.Once
	logger *otelzap.Logger
)

func initLogger() {
	debug := false
	debugFlag := os.Getenv("DEBUG")

	if debugFlag == "true" {
		debug = true
	}

	var (
		log *zap.Logger
		err error
	)

	if debug {
		log, err = zap.NewProduction(zap.AddStacktrace(zap.ErrorLevel), zap.WithCaller(true))
		if err != nil {
			panic(err)
		}
	} else {
		log, err = zap.NewDevelopment(zap.AddStacktrace(zap.ErrorLevel), zap.WithCaller(true))
		if err != nil {
			panic(err)
		}
	}

	logger = otelzap.New(log)
}

func Instance() *otelzap.Logger {
	once.Do(func() {
		initLogger()
	})

	return logger
}
