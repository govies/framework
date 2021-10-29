package config

import (
	"github.com/rs/zerolog"
	"strings"
	"time"
)

type AppConf struct {
	Server  serverConf `yaml:"server"`
	Logging logging    `yaml:"logging"`
}

type logging struct {
	Level string `yaml:"level"`
}

type serverConf struct {
	Mode    string            `yaml:"mode"`
	Port    string            `yaml:"port"`
	Timeout serverTimeoutConf `yaml:"timeout"`
}

type serverTimeoutConf struct {
	Read  time.Duration `yaml:"read"`
	Write time.Duration `yaml:"write"`
	Idle  time.Duration `yaml:"idle"`
}

func (l logging) ZerologLevel() zerolog.Level {
	switch strings.ToLower(l.Level) {
	case "panic":
		return zerolog.PanicLevel
	case "fatal":
		return zerolog.FatalLevel
	case "error":
		return zerolog.ErrorLevel
	case "warn":
		return zerolog.WarnLevel
	case "info":
		return zerolog.InfoLevel
	case "debug":
		return zerolog.DebugLevel
	default:
		return zerolog.TraceLevel
	}
}
