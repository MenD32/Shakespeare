// TraceLogRequest represents a single request from an HTTP trace.
package trace

import (
	"time"
)

// TraceLogRequest represents a single request from an HTTP trace
type TraceLogRequest struct {
	// Timestamp of the request start, relative to the start of the traceLog.
	Delay time.Duration `json:"delay"`

	// HTTP method
	Method string `json:"method"`
	// Path in the URL request, which is the specific endpoint being accessed
	Path string `json:"path"`

	// Request headers
	Headers map[string]string `json:"headers"`

	// Request body
	Body []byte `json:"body"`
}

type TraceLog []TraceLogRequest
