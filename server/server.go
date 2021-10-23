package server

import (
	"github.com/gin-gonic/gin"
	"github.com/govies/framework/config"
	"github.com/govies/framework/logger"
	requestlog "github.com/govies/framework/reguestlog"
	"time"
)

func ListenAndServe() {
	appConf := config.AppConfig()
	log := logger.New(appConf.Logging.ZerologLevel())

	r := gin.New()
	r.Use(requestlog.Logger(log))

	r.GET("/ping", func(c *gin.Context) {
		time.Sleep(time.Second)
		c.JSON(200, gin.H{
			"test": "Hello world!",
		})
	})

	if err := r.Run(":" + appConf.Server.Port); err != nil {
		log.Fatal().Err(err).Msgf("Server startup failed.")
	}

	//mux := http.NewServeMux()
	//mux.HandleFunc("/", Greet)
	//handler := requestlog.NewHandler(mux, log)

	//address := fmt.Sprintf(":%s", appConf.Server.Port)
	//log.Info().Msgf("Starting server on port: %s", appConf.Server.Port)
	//s := &http.Server{
	//	Addr:         address,
	//	Handler:      handler,
	//	ReadTimeout:  appConf.Server.Timeout.Read,
	//	WriteTimeout: appConf.Server.Timeout.Write,
	//	IdleTimeout:  appConf.Server.Timeout.Idle,
	//}
	//if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	//	log.Fatal().Msg("Server startup failed")
	//}
}
