package bootstrap

import (
	"context"
	"os"

	"github.com/gfelixc/9-digits-number/logger"
	"github.com/gfelixc/9-digits-number/server"
)

const LogFilename = "numbers.log"

func StartService(ctx context.Context) error {
	logFile, err := os.Create(LogFilename)
	if err != nil {
		return err
	}

	instrumentedLogger := logger.NewLoggerInstrumented(logger.New(logFile))
	defer instrumentedLogger.Shutdown()

	handler := server.NewReadAndLogLines(instrumentedLogger)
	server := server.New(handler.HandleConnection)

	return server.Start(ctx)
}
