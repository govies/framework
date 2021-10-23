package requestlog

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/govies/framework/logger"
	"net"
	"os"
	"time"
)

func Logger(logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		receivedTime := time.Now()
		// before request
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		hostname, _ := os.Hostname()
		trackingId := c.Request.Header.Get("X-Request-ID")
		if trackingId == "" {
			newUUID, _ := uuid.NewUUID()
			trackingId = newUUID.String()
			c.Request.Header.Set("X-Request-ID", trackingId)
		}

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		requestHeaderMarshal, _ := json.Marshal(c.Request.Header)
		requestBodyMarshal, _ := json.Marshal(c.Request.Body)
		requestLog := logEntry{
			Type:          Request,
			TrackingId:    trackingId,
			ReceivedTime:  receivedTime,
			RequestMethod: c.Request.Method,
			RequestURL:    path,
			RequestBody:   string(requestBodyMarshal),
			RequestHeader: string(requestHeaderMarshal),
			RemoteIP:      ipFromHostPort(c.Request.RemoteAddr),
			HostName:      hostname,
		}
		logSwitch(&requestLog, logger)

		c.Next()
		// after request
		finishedTime := time.Now()
		latency := time.Since(receivedTime)
		// clientIP := c.ClientIP()
		if raw != "" {
			path = path + "?" + raw
		}
		msg := c.Errors.String()
		if msg == "" {
			msg = "Request"
		}

		responseHeaderMarshal, _ := json.Marshal(c.Writer.Header())
		responseLog := logEntry{
			Type:           Response,
			TrackingId:     trackingId,
			FinishedTime:   finishedTime,
			RequestMethod:  c.Request.Method,
			RequestURL:     path,
			Status:         c.Writer.Status(),
			ResponseBody:   blw.body.String(),
			ResponseHeader: string(responseHeaderMarshal),
			RemoteIP:       ipFromHostPort(c.Request.RemoteAddr),
			HostName:       hostname,
		}
		logSwitch(&responseLog, logger)

		summaryLog := logEntry{
			Type:          Summary,
			TrackingId:    trackingId,
			ReceivedTime:  receivedTime,
			FinishedTime:  finishedTime,
			RequestMethod: c.Request.Method,
			RequestURL:    path,
			UserAgent:     c.Request.UserAgent(),
			RemoteIP:      ipFromHostPort(c.Request.RemoteAddr),
			HostName:      hostname,
			Status:        c.Writer.Status(),
			Latency:       latency,
		}
		logSwitch(&summaryLog, logger)
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func logSwitch(data *logEntry, logger *logger.Logger) {
	switch {
	case data.Status >= 400 && data.Status < 500:
		{
			logger.Warn().
				Str("type", string(data.Type)).
				Str("tracking_id", data.TrackingId).
				Time("received_time", data.ReceivedTime).
				Dur("latency", data.Latency).
				Time("finishedTime", data.FinishedTime).
				Int("status", data.Status).
				Str("method", data.RequestMethod).
				Str("url", data.RequestURL).
				Str("response_header", data.ResponseHeader).
				Str("response_body", data.ResponseBody).
				Str("request_header", data.RequestHeader).
				Str("request_body", data.RequestBody).
				Str("remote_ip", data.RemoteIP).
				Str("host_name", data.HostName).
				Str("agent", data.UserAgent).
				Msg("")
		}
	case data.Status >= 500:
		{
			logger.Warn().
				Str("type", string(data.Type)).
				Str("tracking_id", data.TrackingId).
				Time("received_time", data.ReceivedTime).
				Dur("latency", data.Latency).
				Time("finishedTime", data.FinishedTime).
				Int("status", data.Status).
				Str("method", data.RequestMethod).
				Str("url", data.RequestURL).
				Str("response_header", data.ResponseHeader).
				Str("response_body", data.ResponseBody).
				Str("request_header", data.RequestHeader).
				Str("request_body", data.RequestBody).
				Str("remote_ip", data.RemoteIP).
				Str("host_name", data.HostName).
				Str("agent", data.UserAgent).
				Msg("")
		}
	default:
		logger.Info().
			Str("type", string(data.Type)).
			Str("tracking_id", data.TrackingId).
			Time("received_time", data.ReceivedTime).
			Dur("latency", data.Latency).
			Time("finishedTime", data.FinishedTime).
			Int("status", data.Status).
			Str("method", data.RequestMethod).
			Str("url", data.RequestURL).
			Str("response_header", data.ResponseHeader).
			Str("response_body", data.ResponseBody).
			Str("request_header", data.RequestHeader).
			Str("request_body", data.RequestBody).
			Str("remote_ip", data.RemoteIP).
			Str("host_name", data.HostName).
			Str("agent", data.UserAgent).
			Msg("")
	}
}

type RequestResponseType string

const (
	Request  RequestResponseType = "request"
	Response                     = "response"
	Summary                      = "summary"
)

type logEntry struct {
	Type         RequestResponseType
	TrackingId   string
	ReceivedTime time.Time
	Latency      time.Duration
	FinishedTime time.Time

	Status        int
	RequestMethod string
	RequestURL    string

	ResponseBody   string
	ResponseHeader string
	RequestBody    string
	RequestHeader  string

	RemoteIP  string
	HostName  string
	UserAgent string
}

func ipFromHostPort(hp string) string {
	h, _, err := net.SplitHostPort(hp)
	if err != nil {
		return ""
	}
	if len(h) > 0 && h[0] == '[' {
		return h[1 : len(h)-1]
	}
	return h
}
