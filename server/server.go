package server

import (
	"github.com/gin-gonic/gin"
	"github.com/govies/framework/config"
	"github.com/govies/framework/logger"
	"github.com/govies/framework/middlewares"
	"net/http"
)

type Server struct {
	Engine *gin.Engine
}

func New(conf *config.AppConf, l *logger.Logger) *Server {
	gin.SetMode(conf.Server.Mode)
	e := gin.New()

	l.Info().Msg("initializing middlewares")
	e.Use(
		middlewares.Recovery(l, conf),
		middlewares.RequestLogging(l),
	)

	return &Server{Engine: e}
}

func (s *Server) Run(conf *config.AppConf, l *logger.Logger) {
	l.Info().Msgf("server starting on port: %s", conf.Server.Port)
	if err := s.Engine.Run(":" + conf.Server.Port); err != nil {
		l.Fatal().Err(err).Msg("Error while starting server.")
	}
}

func (s *Server) Group(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	return s.Engine.Group(relativePath, handlers...)
}

func (s *Server) Use(middleware ...gin.HandlerFunc) gin.IRoutes {
	return s.Engine.Use(middleware...)
}

func (s *Server) Handle(httpMethod, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.Engine.Handle(httpMethod, relativePath, handlers...)
}
func (s *Server) Any(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.Engine.Any(relativePath, handlers...)
}
func (s *Server) GET(p string, h gin.HandlerFunc) gin.IRoutes {
	return s.Engine.GET(p, h)
}
func (s *Server) POST(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.Engine.POST(relativePath, handlers...)
}
func (s *Server) DELETE(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.Engine.DELETE(relativePath, handlers...)
}
func (s *Server) PATCH(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.Engine.PATCH(relativePath, handlers...)
}
func (s *Server) PUT(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.Engine.PUT(relativePath, handlers...)
}
func (s *Server) OPTIONS(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.Engine.OPTIONS(relativePath, handlers...)
}
func (s *Server) HEAD(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.Engine.HEAD(relativePath, handlers...)
}
func (s *Server) StaticFile(relativePath, filepath string) gin.IRoutes {
	return s.Engine.StaticFile(relativePath, filepath)
}
func (s *Server) Static(relativePath, root string) gin.IRoutes {
	return s.Engine.Static(relativePath, root)
}
func (s *Server) StaticFS(relativePath string, fs http.FileSystem) gin.IRoutes {
	return s.Engine.StaticFS(relativePath, fs)
}
