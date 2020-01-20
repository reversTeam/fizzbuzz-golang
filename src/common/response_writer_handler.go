package common

import (
	"net/http"
)

// Definition of ResponseWriterHandler
type ResponseWriterHandler struct {
    http.ResponseWriter
    StatusCode int
}

// Create a ResponseWriterHandler
func NewResponseWriterHandler(w http.ResponseWriter) *ResponseWriterHandler {
    return &ResponseWriterHandler{w, http.StatusOK}
}

// Required by the Exporter for catch status code of Grpc request
func (o *ResponseWriterHandler) WriteHeader(code int) {
    o.StatusCode = code
    o.ResponseWriter.WriteHeader(code)
}