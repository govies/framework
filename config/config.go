package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

func DefaultConf() *AppConf {
	return &AppConf{
		Server: serverConf{
			Mode: "release",
			Port: "8080",
			Timeout: serverTimeoutConf{
				Read:  30 * time.Second,
				Write: 30 * time.Second,
				Idle:  120 * time.Second,
			},
		},
		Logging: logging{
			Level: "info",
		},
	}
}

func DefaultAppConf() *AppConf {
	conf := DefaultConf()
	conf.Load(DefaultFile())
	return conf
}

type File struct {
	Path string
	Name string
	Type string
}

func DefaultFile() *File {
	return &File{
		Path: ".",
		Name: "configs",
		Type: "yaml",
	}
}

type Config interface {
	Load(f *File)
}

func (c *AppConf) Load(f *File) {
	LoadConfigFile(c, f)
}

func LoadConfigFile(c interface{}, f *File) {
	viper.AddConfigPath(f.Path)
	viper.SetConfigName(f.Name)
	viper.SetConfigType(f.Type)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("%v", err)
	}
	if err := viper.Unmarshal(c); err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}
}
