// TraceLogRequest represents a single request from an HTTP trace.
package clients

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
	"time"
)

type Trace interface {
	Dump(*os.File) error
}

// TraceLogRequest represents a single request from an HTTP trace
type TraceLogRequest struct {
	// Timestamp of the request start, relative to the start of the traceLog.
	Delay time.Duration `json:"delay"`

	// HTTP method
	Method string `json:"method"`
	// Path in the URL request, which is the specific endpoint being accessed
	Path string `json:"path"`

	// Protocol specification
	Proto      string `json:"proto"` // e.g. "HTTP/1.1"
	ProtoMajor int    `json:"major"` // e.g. 1
	ProtoMinor int    `json:"minor"` // e.g. 1

	// Request headers
	Header []http.Header `json:"headers"`

	// Request body
	Body []byte `json:"body"`
}

type TraceLog struct {
	Requests []TraceLogRequest
}

func New(filepath string) (*TraceLog, error) {
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

func (t *TraceLog) Dump(file *os.File) error {
	for _, request := range t.Requests {
		data, err := json.Marshal(request)
		if err != nil {
			return err
		}
		_, err = file.Write(append(data, '\n'))
		if err != nil {
			return err
		}
	}
	return nil
}
