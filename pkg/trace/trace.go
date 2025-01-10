// TraceLogRequest represents a single request from an HTTP trace.
package trace

import (
	"net/http"
	"time"
)

// TraceLogRequest represents a single request from an HTTP trace
type TraceLogRequest struct {
	// Timestamp of the request start, relative to the start of the traceLog.
	Timestamp time.Duration

	// HTTP method
	Method string
	// Path in the URL request, which is the specific endpoint being accessed
	Path string

	// Protocol specification
	Proto      string // e.g. "HTTP/1.1"
	ProtoMajor int    // e.g. 1
	ProtoMinor int    // e.g. 1

	// Request headers
	Header http.Header

	// Request body
	Body []byte
}

type TraceLog struct {
	Requests []TraceLogRequest
}
