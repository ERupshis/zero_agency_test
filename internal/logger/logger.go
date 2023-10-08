package logger

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// loggerZap wrapper of Zap logger.
type loggerZap struct {
	zap *zap.Logger
}

// CreateZapLogger create method for zap logger.
func CreateZapLogger(level string) (BaseLogger, error) {
	cfg, err := initConfig(level)
	if err != nil {
		return nil, err
	}

	log, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("create zap loggerZap^ %w", err)
	}

	return &loggerZap{zap: log}, nil
}

// Info generates 'info' level log.
func (l *loggerZap) Info(msg string, fields ...interface{}) {
	l.zap.Info(fmt.Sprintf(msg, fields...))
}

// Printf interface for kafka's implementation.
func (l *loggerZap) Printf(msg string, fields ...interface{}) {
	l.Info(msg, fields...)
}

// initConfig method that initializes logger.
func initConfig(level string) (zap.Config, error) {
	cfg := zap.NewProductionConfig()

	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		emptyConfig := zap.Config{}
		return emptyConfig, fmt.Errorf("init zap loggerZap config: %w", err)
	}
	cfg.Level = lvl
	cfg.DisableCaller = true

	return cfg, nil
}

// Sync flushing any buffered log entries.
func (l *loggerZap) Sync() {
	err := l.zap.Sync()
	if err != nil {
		panic(err)
	}
}

// LogHandler handler for requests logging.
func (l *loggerZap) LogHandler(c *fiber.Ctx) error {
	start := time.Now()

	// Proceed with the request handling
	err := c.Next()

	// Calculate the duration
	duration := time.Since(start)

	// Log the request details
	l.zap.Info("new incoming HTTP request",
		zap.String("uri", c.OriginalURL()),
		zap.String("method", c.Method()),
		zap.Int("status", c.Response().StatusCode()), // Capture the status code from the response
		zap.String("content-type", c.GetRespHeader("Content-Type")),
		zap.Duration("duration", duration),
		zap.Int("size", len(c.Response().Body())),
	)

	return err
}
