package server

import (
	"github.com/gin-gonic/gin"
	"github.com/govies/framework/config"
	"github.com/govies/framework/handler"
	"github.com/govies/framework/logger"
	"net/http"
)

type Server struct {
	router *gin.Engine
}

var (
	configs = config.AppConfig()
	log     = logger.New()
)

func NewServer() *Server {
	gin.SetMode(config.AppConfig().Logging.Level)
	e := gin.New()
	e.Use(handler.RequestHandler())
	e.Use(handler.ErrorHandler())
	e.Use(handler.RecoveryHandler())
	return &Server{
		router: e,
	}
}

func (s *Server) Run() error {
	log.Info().Msgf("Server starting in %s mode on port: %s.", configs.Server.Mode, configs.Server.Port)
	err := s.router.Run(":" + configs.Server.Port)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Group(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	return s.router.Group(relativePath, handlers...)
}

func (s *Server) Use(middleware ...gin.HandlerFunc) gin.IRoutes {
	return s.router.Use(middleware...)
}

func (s *Server) Handle(httpMethod, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.router.Handle(httpMethod, relativePath, handlers...)
}
func (s *Server) Any(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.router.Any(relativePath, handlers...)
}
func (s *Server) GET(p string, h gin.HandlerFunc) gin.IRoutes {
	return s.router.GET(p, h)
}
func (s *Server) POST(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.router.POST(relativePath, handlers...)
}
func (s *Server) DELETE(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.router.DELETE(relativePath, handlers...)
}
func (s *Server) PATCH(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.router.PATCH(relativePath, handlers...)
}
func (s *Server) PUT(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.router.PUT(relativePath, handlers...)
}
func (s *Server) OPTIONS(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.router.OPTIONS(relativePath, handlers...)
}
func (s *Server) HEAD(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.router.HEAD(relativePath, handlers...)
}
func (s *Server) StaticFile(relativePath, filepath string) gin.IRoutes {
	return s.router.StaticFile(relativePath, filepath)
}
func (s *Server) Static(relativePath, root string) gin.IRoutes {
	return s.router.Static(relativePath, root)
}
func (s *Server) StaticFS(relativePath string, fs http.FileSystem) gin.IRoutes {
	return s.router.StaticFS(relativePath, fs)
}
