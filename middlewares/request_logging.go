package middlewares

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/govies/onboard/logger"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

func RequestLogging(l *logger.Logger) gin.HandlerFunc {
	l.Info().Msg("initializing request logging middleware")
	return func(c *gin.Context) {
		receivedTime := time.Now()
		hostname, _ := os.Hostname()
		baseLogEntry := &logger.BaseLogEntry{
			TrackingId:    GetTrackingId(&c.Request.Header),
			RequestURL:    GetFullPath(c.Request.URL),
			UserAgent:     c.Request.UserAgent(),
			HostName:      hostname,
			RemoteIP:      IpFromHostPort(c.Request.RemoteAddr),
			RequestMethod: c.Request.Method,
		}

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		requestHeaderMarshal, _ := json.Marshal(c.Request.Header)
		requestBodyMarshal, _ := json.Marshal(c.Request.Body)
		baseLogEntry.Type = logger.Request
		requestLog := &logger.requestLogEntry{
			BaseLogEntry:  baseLogEntry,
			ReceivedTime:  receivedTime,
			RequestBody:   string(requestBodyMarshal),
			RequestHeader: string(requestHeaderMarshal),
		}
		requestLog.Log(l)

		c.Next()
		latency := time.Since(receivedTime)
		finishedTime := receivedTime.Add(latency)

		responseHeaderMarshal, _ := json.Marshal(c.Writer.Header())
		baseLogEntry.Type = logger.Response
		responseLogEntry := logger.responseLogEntry{
			BaseLogEntry:   baseLogEntry,
			Status:         c.Writer.Status(),
			ResponseBody:   blw.body.String(),
			ResponseHeader: string(responseHeaderMarshal),
			FinishedTime:   finishedTime,
		}
		responseLogEntry.Log(l)

		baseLogEntry.Type = logger.Summary
		summaryLogEntry := logger.summaryLogEntry{
			BaseLogEntry: baseLogEntry,
			Status:       c.Writer.Status(),
			Latency:      latency,
			FinishedTime: finishedTime,
			ReceivedTime: receivedTime,
		}
		summaryLogEntry.Log(l)

		err := c.Errors.Last()
		if err != nil {
			baseLogEntry.Type = logger.Error
			errorLogEntry := logger.ErrorLogEntry{
				BaseLogEntry: baseLogEntry,
				ReceivedTime: receivedTime,
				Error:        err.Unwrap(),
				FinishedTime: finishedTime,
				Latency:      latency,
				Status:       c.Writer.Status(),
			}
			errorLogEntry.Log(l)
		}
	}
}

func GetFullPath(u *url.URL) string {
	path := u.Path
	raw := u.RawQuery
	if raw != "" {
		return path + "?" + raw
	}
	return path
}

func GetTrackingId(h *http.Header) string {
	trackingId := h.Get("X-Request-ID")
	if trackingId == "" {
		newUUID, _ := uuid.NewUUID()
		trackingId = newUUID.String()
		h.Set("X-Request-ID", trackingId)
	}
	return trackingId
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func IpFromHostPort(hp string) string {
	h, _, err := net.SplitHostPort(hp)
	if err != nil {
		return ""
	}
	if len(h) > 0 && h[0] == '[' {
		return h[1 : len(h)-1]
	}
	return h
}
