package log

import (
	"os"

	"go.uber.org/zap/zapcore"
)

type Option func(cfg *config)

type config struct {
	level           Level
	TimeEncoder     TimeEncoder
	LevelEncoder    LevelEncoder
	DurationEncoder DurationEncoder
	CallerEncoder   CallerEncoder
	skipCaller      int
	encoder         string
	w               zapcore.WriteSyncer

	cores []Core
	hooks []Hook
}

var defaultConfig = config{
	level:           InfoLevel,
	skipCaller:      1,
	encoder:         "console",
	TimeEncoder:     timeEncoder,
	LevelEncoder:    zapcore.CapitalLevelEncoder,
	DurationEncoder: zapcore.StringDurationEncoder,
	CallerEncoder:   zapcore.ShortCallerEncoder,
	w:               os.Stdout,
}

func WithLevel(lvl Level) Option {
	return func(cfg *config) {
		cfg.level = lvl
	}
}

func WithSkipCaller(skip int) Option {
	return func(cfg *config) {
		cfg.skipCaller = skip
	}
}

func WithJSON() Option {
	return func(cfg *config) {
		cfg.encoder = "json"
	}
}

func WithSensitive() Option {
	return func(cfg *config) {
		cfg.encoder = "sensitive"
	}
}

func WithConsole() Option {
	return func(cfg *config) {
		cfg.encoder = "console"
	}
}

func WithCore(cores ...Core) Option {
	return func(cfg *config) {
		cfg.cores = append(cfg.cores, cores...)
	}
}

func WithHook(hooks ...Hook) Option {
	return func(cfg *config) {
		cfg.hooks = append(cfg.hooks, hooks...)
	}
}

func WithWriter(w zapcore.WriteSyncer) Option {
	return func(cfg *config) {
		cfg.w = w
	}
}
