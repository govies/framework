package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/govies/framework/logger"
	"github.com/govies/framework/resp"
	"github.com/govies/framework/server"
	"net/http"
)

var (
	log = logger.New()
)

func main() {
	s := server.NewServer()

	v1 := s.Group("/api/v1")
	{
		v1.GET("/hiba", funcName())
	}

	s.Group("/api/v2")
	{
		s.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"hello": "world",
			})
		})
	}

	if err := s.Run(); err != nil {
		log.Fatal().Err(err).Msg("Server failed to start.")
	}
}

func funcName() gin.HandlerFunc {
	return func(c *gin.Context) {
		//response := resp.Success(200, User{Id: "1", Name: "valaki"})
		response := resp.Error(http.StatusForbidden, errors.New("valami"))
		response.Send(c)
		//response := resp.Response{Status: 200, Data: User{Id: "1", Name: "valaki"}, Errors: error.Success("valami hiba")}
		//response.Send(c)
		return
	}
}

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
