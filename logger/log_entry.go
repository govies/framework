package logger

import (
	"github.com/rs/zerolog"
	"time"
)

type RequestResponseType string

const (
	Request  RequestResponseType = "request"
	Response                     = "resp"
	Summary                      = "summary"
	Error                        = "error"
)

type Logging interface {
	Log(l *Logger)
}

type BaseLogEntry struct {
	Type       RequestResponseType
	TrackingId string

	RequestMethod string
	RequestURL    string

	RemoteIP  string
	HostName  string
	UserAgent string
}

type RequestLogEntry struct {
	BaseLogEntry  *BaseLogEntry
	ReceivedTime  time.Time
	RequestBody   string
	RequestHeader string
}

func (r *RequestLogEntry) Log(l *Logger) {
	l.Info().
		Str("type", string(r.BaseLogEntry.Type)).
		Str("tracking_id", r.BaseLogEntry.TrackingId).
		Time("received_time", r.ReceivedTime).
		Str("method", r.BaseLogEntry.RequestMethod).
		Str("url", r.BaseLogEntry.RequestURL).
		Str("request_header", r.RequestHeader).
		Str("request_body", r.RequestBody).
		Str("remote_ip", r.BaseLogEntry.RemoteIP).
		Str("host_name", r.BaseLogEntry.HostName).
		Str("agent", r.BaseLogEntry.UserAgent).
		Msg("")
}

type ResponseLogEntry struct {
	BaseLogEntry   *BaseLogEntry
	FinishedTime   time.Time
	Status         int
	ResponseBody   string
	ResponseHeader string
}

func (r ResponseLogEntry) Log(l *Logger) {
	l.LevelByStatus(r.Status).
		Str("type", string(r.BaseLogEntry.Type)).
		Str("tracking_id", r.BaseLogEntry.TrackingId).
		Time("finished_time", r.FinishedTime).
		Int("status", r.Status).
		Str("method", r.BaseLogEntry.RequestMethod).
		Str("url", r.BaseLogEntry.RequestURL).
		Str("response_header", r.ResponseHeader).
		Str("response_body", r.ResponseBody).
		Str("remote_ip", r.BaseLogEntry.RemoteIP).
		Str("host_name", r.BaseLogEntry.HostName).
		Str("agent", r.BaseLogEntry.UserAgent).
		Msg("")
}

type SummaryLogEntry struct {
	BaseLogEntry *BaseLogEntry
	Status       int
	ReceivedTime time.Time
	FinishedTime time.Time
	Latency      time.Duration
}

func (s SummaryLogEntry) Log(l *Logger) {
	l.LevelByStatus(s.Status).
		Str("type", string(s.BaseLogEntry.Type)).
		Str("tracking_id", s.BaseLogEntry.TrackingId).
		Time("received_time", s.ReceivedTime).
		Time("finished_time", s.FinishedTime).
		Dur("latency", s.Latency).
		Int("status", s.Status).
		Str("method", s.BaseLogEntry.RequestMethod).
		Str("url", s.BaseLogEntry.RequestURL).
		Str("remote_ip", s.BaseLogEntry.RemoteIP).
		Str("host_name", s.BaseLogEntry.HostName).
		Str("agent", s.BaseLogEntry.UserAgent).
		Msg("")
}

type ErrorLogEntry struct {
	BaseLogEntry *BaseLogEntry
	ReceivedTime time.Time
	FinishedTime time.Time
	Latency      time.Duration
	Status       int
	Error        error
	Stack        string
}

func (e ErrorLogEntry) Log(l *Logger) {
	if l.LogLevel <= zerolog.DebugLevel {
		l.LevelByStatus(e.Status).
			Str("type", string(e.BaseLogEntry.Type)).
			Str("tracking_id", e.BaseLogEntry.TrackingId).
			Time("received_time", e.ReceivedTime).
			Time("finished_time", e.FinishedTime).
			Dur("latency", e.Latency).
			Int("status", e.Status).
			Str("method", e.BaseLogEntry.RequestMethod).
			Str("url", e.BaseLogEntry.RequestURL).
			AnErr("error", e.Error).
			Str("stack", e.Stack).
			Str("remote_ip", e.BaseLogEntry.RemoteIP).
			Str("host_name", e.BaseLogEntry.HostName).
			Str("agent", e.BaseLogEntry.UserAgent).
			Msg("")
	} else {
		l.LevelByStatus(e.Status).
			Str("type", string(e.BaseLogEntry.Type)).
			Str("tracking_id", e.BaseLogEntry.TrackingId).
			Time("received_time", e.ReceivedTime).
			Time("finished_time", e.FinishedTime).
			Dur("latency", e.Latency).
			Int("status", e.Status).
			Str("method", e.BaseLogEntry.RequestMethod).
			Str("url", e.BaseLogEntry.RequestURL).
			AnErr("error", e.Error).
			Str("remote_ip", e.BaseLogEntry.RemoteIP).
			Str("host_name", e.BaseLogEntry.HostName).
			Str("agent", e.BaseLogEntry.UserAgent).
			Msg("")
	}
}
