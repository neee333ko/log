package log

import (
	"github.com/marmotedu/component-base/pkg/validation/field"
	"github.com/spf13/pflag"
)

type Options struct {
	Level       string `json:"lvl" mapstructure:"lvl"`
	Development bool   `json:"develop" mapstructure:"develop"`
	// support json or console
	Format            string   `json:"fmt" mapstructure:"fmt"`
	DisableCaller     bool     `json:"disableCaller" mapstructure:"disableCaller"`
	DisableStacktrace bool     `json:"disableStacktrace" mapstructure:"disableStacktrace"`
	EnableColor       bool     `json:"enableColor" mapstructure:"enableColor"`
	Name              string   `json:"name" mapstructure:"name"`
	OutputPaths       []string `json:"paths" mapstructure:"paths"`
	ErrorOutputPaths  []string `json:"errpaths" mapstructure:"errpaths"`
}

type operation func(*Options)

func InitOptions(ops ...operation) *Options {
	opts := &Options{
		Level:             "info",
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

func WithLevel(Level string) operation {
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

func (o *Options) Flags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("LogOption", pflag.ExitOnError)

	fs.StringVar(&o.Level, "log.lvl", "info", "logger level")
	fs.BoolVar(&o.Development, "log.develop", false, "whether set logger in development mode")
	fs.StringVar(&o.Format, "log.fmt", "json", "output format, support json and console")
	fs.BoolVar(&o.DisableCaller, "log.disableCaller", false, "disable output logger caller")
	fs.BoolVar(&o.DisableStacktrace, "log.disableStacktrace", false, "disable output logger stacktrace")
	fs.BoolVar(&o.EnableColor, "log.enableColor", true, "enable color output")
	fs.StringVar(&o.Name, "log.name", "log", "logger name")
	fs.StringSliceVar(&o.OutputPaths, "log.paths", []string{"stdout"}, "logger output paths")
	fs.StringSliceVar(&o.ErrorOutputPaths, "log.errpaths", []string{"stderr"}, "logger err output paths")

	return fs
}

func (o *Options) Validate() field.ErrorList {
	return nil
}

func (o *Options) Complete() error {
	return nil
}
