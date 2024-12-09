package utils

import (
	"bytes"
	"net"
	"net/http"
	"path/filepath"
	"strings"
)

// ClientIP get client ip
func ClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

// ResponseCapture capture response and get response info
type ResponseCapture struct {
	http.ResponseWriter
	wroteHeader bool
	status      int
	body        *bytes.Buffer
}

// NewResponseCapture new response capture
func NewResponseCapture(w http.ResponseWriter) *ResponseCapture {
	return &ResponseCapture{
		ResponseWriter: w,
		wroteHeader:    false,
		body:           new(bytes.Buffer),
	}
}

// Header return response writer header
func (c ResponseCapture) Header() http.Header {
	return c.ResponseWriter.Header()
}

// Write write data to response writer and body
func (c ResponseCapture) Write(data []byte) (int, error) {
	if !c.wroteHeader {
		c.WriteHeader(http.StatusOK)
	}
	c.body.Write(data)
	return c.ResponseWriter.Write(data)
}

// WriteHeader write header to response writer
func (c *ResponseCapture) WriteHeader(statusCode int) {
	c.status = statusCode
	c.wroteHeader = true
	c.ResponseWriter.WriteHeader(statusCode)
}

// Bytes return response body bytes
func (c ResponseCapture) Bytes() []byte {
	return c.body.Bytes()
}

// StatusCode return status code
func (c ResponseCapture) StatusCode() int {
	return c.status
}

// CleanRelativePath returns the shortest path name equivalent to path
// by purely lexical processing. It make sure the provided path is rooted
// and then uses filepath.Clean and filepath.Rel to make sure the path
// doesn't include any separators or elements that shouldn't be there
// like ., .., //.
func CleanRelativePath(path string) (string, error) {
	cleanPath := filepath.Clean(filepath.Join("/", path))
	rel, err := filepath.Rel("/", cleanPath)
	if err != nil {
		// slash is prepended above therefore this is not expected to fail
		return "", err
	}

	return rel, nil
}
