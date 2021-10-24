package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

func AppConfig() *Conf {
	viper.AddConfigPath(".")
	viper.SetConfigName("configs")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("%v", err)
	}
	conf := defaultConf()
	if err := viper.Unmarshal(conf); err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}
	return conf
}

func defaultConf() *Conf {
	return &Conf{
		Server: serverConf{
			Mode: "release",
			Port: "8080",
			Timeout: serverTimeoutConf{
				Read:  30 * time.Second,
				Write: 30 * time.Second,
				Idle:  120 * time.Second,
			},
		},
	}
}
