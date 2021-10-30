package server

import (
	"github.com/gin-gonic/gin"
	"github.com/govies/framework/config"
	"github.com/govies/framework/logger"
	"github.com/govies/framework/middlewares"
)

type Server struct {
	Engine *gin.Engine
	Conf   *config.AppConf
	Log    *logger.Logger
}

func New(conf *config.AppConf, l *logger.Logger) *Server {
	gin.SetMode(conf.Server.Mode)
	e := gin.New()
	return &Server{
		Engine: e,
		Conf:   conf,
		Log:    l,
	}
}

func (s *Server) Run() {
	s.Log.Info().Msgf("server starting on port: %s", s.Conf.Server.Port)
	if err := s.Engine.Run(":" + s.Conf.Server.Port); err != nil {
		s.Log.Fatal().Err(err).Msg("Error while starting server.")
	}
}

func (s *Server) InitDefaultMiddlewares() {
	s.Log.Info().Msg("initializing middlewares")
	s.Engine.Use(
		middlewares.Recovery(s.Log, s.Conf),
		middlewares.RequestLogging(s.Log),
	)
}
