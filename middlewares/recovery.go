package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"github.com/govies/framework/config"
	"github.com/govies/framework/error"
	"github.com/govies/framework/logger"
	"github.com/govies/framework/resp"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"time"
)

func Recovery(l *logger.Logger, cfg *config.AppConf) gin.HandlerFunc {
	l.Info().Msg("initializing recovery middleware")
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				goErr := errors.Wrap(err, 3)
				if l != nil {
					logError(c, l, goErr)
				} else {
					fmt.Printf("Received: %v, ErrorDto message: %v, Stack: %v", time.Now(), goErr, goErr.ErrorStack())
				}
				sendServerError(c, cfg, l, goErr)
			}
		}()
		c.Next()
	}
}

func sendServerError(c *gin.Context, cfg *config.AppConf, l *logger.Logger, goErr *errors.Error) {
	errorDto := error.FromErrors(http.StatusInternalServerError, goErr)
	errorDto.UserMessage = "Something went wrong"

	if cfg.Logging.ZerologLevel() <= zerolog.DebugLevel {
		errorDto.Stack = goErr.ErrorStack()
		resp.ErrorDto(http.StatusInternalServerError, errorDto).Send(c, l)
	} else {
		resp.ErrorDto(http.StatusInternalServerError, errorDto).Send(c, l)
	}
}

func logError(c *gin.Context, l *logger.Logger, goErr *errors.Error) {
	hostname, _ := os.Hostname()

	baseLogEntry := &logger.BaseLogEntry{
		Type:          "error",
		TrackingId:    GetTrackingId(&c.Request.Header),
		RequestURL:    GetFullPath(c.Request.URL),
		UserAgent:     c.Request.UserAgent(),
		HostName:      hostname,
		RemoteIP:      IpFromHostPort(c.Request.RemoteAddr),
		RequestMethod: c.Request.Method,
	}

	errorLogEntry := logger.ErrorLogEntry{
		BaseLogEntry: baseLogEntry,
		Error:        goErr,
		Status:       500,
		Stack:        goErr.ErrorStack(),
		ReceivedTime: time.Now(),
	}
	errorLogEntry.Log(l)
}
