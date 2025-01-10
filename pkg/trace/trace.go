// TraceLogRequest represents a single request from an HTTP trace.
package trace

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
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

func NewTraceLog(filepath string) (*TraceLog, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var requests []TraceLogRequest
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		request := TraceLogRequest{}
		err = json.Unmarshal([]byte(line), &request)
		if err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &TraceLog{Requests: requests}, nil
}
