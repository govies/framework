package server

import (
	"github.com/govies/framework/config"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
)

func ListenAndServe() {
	setupViperConfig()
	startServer()
}

func startServer() {
	serverPort := config.GetStringOrDefault("server.port", "8080")
	server := http.Server{
		Addr:    ":" + serverPort,
		Handler: nil,
	}
	log.Println("Server started on port: ", serverPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func setupViperConfig() {
	viper.SetConfigType("yaml")
	if configPath := os.Getenv("CONFIG_FILE_PATH"); configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		viper.SetConfigFile("configs.yaml")
	}
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Cannot read config file. Error: ", err)
	}
}
