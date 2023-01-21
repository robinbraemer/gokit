package log

import (
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config are the options for creating a logger with New.
type Config struct {
	Debug     bool
	Verbosity int
}

// OptionFn is a function that configures a logger for New.
type OptionFn func(*Config)

// WithConfig overrides all options.
func WithConfig(config *Config) OptionFn {
	return func(c *Config) {
		*c = *config
	}
}

// New creates a new logger with the given options.
func New(opts ...OptionFn) (l logr.Logger, err error) {
	cfg := new(Config)
	for _, o := range opts {
		o(cfg)
	}

	var zapCfg zap.Config
	if cfg.Debug { // always use debug style logs for now
		zapCfg = zap.NewDevelopmentConfig()
		zapCfg.Encoding = "console"
		zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zapCfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000") // zapcore.ISO8601TimeEncoder
	} else {
		zapCfg = zap.NewProductionConfig()
	}
	zapCfg.Level = zap.NewAtomicLevelAt(zapcore.Level(-cfg.Verbosity))

	zl, err := zapCfg.Build()
	if err != nil {
		return logr.Discard(), err
	}
	return zapr.NewLogger(zl), nil
}
