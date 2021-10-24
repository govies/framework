package error

import (
	"net/http"
	"time"
)

type Dto struct {
	Status        int      `json:"status,omitempty"`
	Code          string   `json:"code,omitempty"`
	Timestamp     int64    `json:"timestamp,omitempty"`
	DebugMessages []string `json:"DebugMessages,omitempty"`
	UserMessage   string   `json:"userMessage,omitempty"`
	Stack         string   `json:"stack,omitempty"`
}

func FromErrors(s int, errors ...error) *Dto {
	dto := New(s)
	for _, v := range errors {
		dto.AppendDebugMessages(v.Error())
	}
	return dto
}

func New(s int) *Dto {
	return &Dto{
		Status:      s,
		Code:        http.StatusText(s),
		Timestamp:   time.Now().UnixMilli(),
		UserMessage: "Something went wrong.",
	}
}

func (e Dto) AppendDebugMessages(m string) *Dto {
	e.DebugMessages = append(e.DebugMessages, m)
	return &e
}
