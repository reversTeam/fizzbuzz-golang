package common

import (
	"net/http"
)

type ResponseWriterHandler struct {
    http.ResponseWriter
    StatusCode int
}

func NewResponseWriterHandler(w http.ResponseWriter) *ResponseWriterHandler {
    return &ResponseWriterHandler{w, http.StatusOK}
}

func (o *ResponseWriterHandler) WriteHeader(code int) {
    o.StatusCode = code
    o.ResponseWriter.WriteHeader(code)
}