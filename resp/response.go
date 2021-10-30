package resp

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	err "github.com/govies/framework/error"
)

type Response struct {
	httpStatus int
	Success    bool        `json:"success,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Errors     *err.Dto    `json:"error,omitempty"`
}

func Success(s int, d interface{}) *Response {
	return &Response{
		httpStatus: s,
		Success:    true,
		Data:       d,
		Errors:     nil,
	}
}

func Error(s int, e ...error) *Response {
	return &Response{
		httpStatus: s,
		Success:    false,
		Errors:     err.FromErrors(s, e...),
	}
}

func ErrorDto(s int, e *err.Dto) *Response {
	return &Response{
		httpStatus: s,
		Success:    false,
		Errors:     e,
	}
}

func (r *Response) Send(c *gin.Context) {
	if !r.Success {
		marshal, _ := json.Marshal(r.Errors)
		_ = c.Error(errors.New(string(marshal)))
	}
	c.JSON(r.httpStatus, r)
}
