// TraceLogRequest represents a single request from an HTTP trace.
package trace

import (
	"bufio"
	"encoding/json"
	"os"
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

func Load(filepath string) (TraceLog, error) {
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
	return requests, nil
}

func (t TraceLog) Dump(file *os.File) error {
	for _, request := range t {
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
