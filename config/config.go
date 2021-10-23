package config

import "github.com/spf13/viper"

func GetStringOrDefault(k string, d string) string {
	if v := viper.GetString(k); v != "" {
		return v
	}
	return d
}
