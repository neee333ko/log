package log

import (
	"go.uber.org/zap/zapcore"
)

type Options struct {
	Level             zapcore.Level
	Development       bool
	Format            string
	DisableCaller     bool
	DisableStacktrace bool
	EnableColor       bool
	Name              string
	OutputPaths       []string
	ErrorOutputPaths  []string
}

type operation func(*Options)

func InitOptions(ops ...operation) *Options {
	opts := &Options{
		Level:             zapcore.InfoLevel,
		Development:       false,
		Format:            "json",
		DisableCaller:     false,
		DisableStacktrace: false,
		EnableColor:       true,
		Name:              "Logger",
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
	}

	for _, op := range ops {
		op(opts)
	}

	return opts
}

func WithLevel(Level zapcore.Level) operation {
	return func(o *Options) {
		o.Level = Level
	}
}

func WithDevelopment(Development bool) operation {
	return func(o *Options) {
		o.Development = Development
	}
}

func WithFormat(Format string) operation {
	return func(o *Options) {
		o.Format = Format
	}
}

func WithDisableCaller(DisableCaller bool) operation {
	return func(o *Options) {
		o.DisableCaller = DisableCaller
	}
}

func WithDisableStacktrace(DisableStacktrace bool) operation {
	return func(o *Options) {
		o.DisableStacktrace = DisableStacktrace
	}
}

func WithEnableColor(EnableColor bool) operation {
	return func(o *Options) {
		o.EnableColor = EnableColor
	}
}

func WithName(Name string) operation {
	return func(o *Options) {
		o.Name = Name
	}
}

func WithOutputPaths(OutputPaths []string) operation {
	return func(o *Options) {
		o.OutputPaths = OutputPaths
	}
}

func WithErrorOutputPaths(ErrorOutputPaths []string) operation {
	return func(o *Options) {
		o.ErrorOutputPaths = ErrorOutputPaths
	}
}
