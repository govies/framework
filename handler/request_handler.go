package handler

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

func RequestHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		receivedTime := time.Now()
		hostname, _ := os.Hostname()
		baseLogEntry := &BaseLogEntry{
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
		baseLogEntry.Type = Request
		requestLog := &requestLogEntry{
			BaseLogEntry:  baseLogEntry,
			ReceivedTime:  receivedTime,
			RequestBody:   string(requestBodyMarshal),
			RequestHeader: string(requestHeaderMarshal),
		}
		requestLog.log()

		c.Next()
		latency := time.Since(receivedTime)
		finishedTime := receivedTime.Add(latency)

		responseHeaderMarshal, _ := json.Marshal(c.Writer.Header())
		baseLogEntry.Type = Response
		responseLogEntry := responseLogEntry{
			BaseLogEntry:   baseLogEntry,
			Status:         c.Writer.Status(),
			ResponseBody:   blw.body.String(),
			ResponseHeader: string(responseHeaderMarshal),
			FinishedTime:   finishedTime,
		}
		responseLogEntry.log()

		baseLogEntry.Type = Summary
		summaryLogEntry := summaryLogEntry{
			BaseLogEntry: baseLogEntry,
			Status:       c.Writer.Status(),
			Latency:      latency,
			FinishedTime: finishedTime,
			ReceivedTime: receivedTime,
		}
		summaryLogEntry.log()
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
