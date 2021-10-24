package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"github.com/govies/framework/error"
	"github.com/rs/zerolog"
	"os"
	"time"
)

func RecoveryHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				goErr := errors.Wrap(err, 3)
				if log != nil {

					hostname, _ := os.Hostname()

					baseLogEntry := &BaseLogEntry{
						Type:          "error",
						TrackingId:    GetTrackingId(&c.Request.Header),
						RequestURL:    GetFullPath(c.Request.URL),
						UserAgent:     c.Request.UserAgent(),
						HostName:      hostname,
						RemoteIP:      IpFromHostPort(c.Request.RemoteAddr),
						RequestMethod: c.Request.Method,
					}

					errorLogEntry := ErrorLogEntry{
						BaseLogEntry: baseLogEntry,
						Error:        goErr,
						Status:       500,
						Stack:        goErr.ErrorStack(),
						ReceivedTime: time.Now(),
					}
					errorLogEntry.Log()
				} else {
					fmt.Printf("Received: %v, ErrorDto message: %v, Stack: %v", time.Now(), goErr, goErr.ErrorStack())
				}

				errorDto := error.FromErrors(500, goErr)
				errorDto.UserMessage = "Something went wrong"

				if conf.Logging.ZerologLevel() <= zerolog.DebugLevel {
					errorDto.Stack = goErr.ErrorStack()
					c.JSON(500, errorDto)
				} else {
					c.JSON(500, errorDto)
				}
			}
		}()
		c.Next() // execute all the handlers
	}
}
