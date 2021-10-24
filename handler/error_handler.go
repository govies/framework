package handler

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"os"
	"time"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		receivedTime := time.Now()
		hostname, _ := os.Hostname()

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()
		latency := time.Since(receivedTime)
		finishedTime := receivedTime.Add(latency)
		err := c.Errors.Last()
		if err == nil {
			return
		}

		baseLogEntry := &BaseLogEntry{
			Type:          Error,
			TrackingId:    GetTrackingId(&c.Request.Header),
			RequestURL:    GetFullPath(c.Request.URL),
			UserAgent:     c.Request.UserAgent(),
			HostName:      hostname,
			RemoteIP:      IpFromHostPort(c.Request.RemoteAddr),
			RequestMethod: c.Request.Method,
		}

		errorLogEntry := ErrorLogEntry{
			BaseLogEntry: baseLogEntry,
			ReceivedTime: receivedTime,
			Error:        err.Unwrap(),
			FinishedTime: finishedTime,
			Latency:      latency,
			Status:       500,
		}
		errorLogEntry.Log()
	}
}
