package runtime

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SSEWriter struct {
	W http.ResponseWriter
}

func (s SSEWriter) Event(name string, data any) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(s.W, "event: %s\ndata: %s\n\n", name, string(bytes)); err != nil {
		return err
	}
	if flusher, ok := s.W.(http.Flusher); ok {
		flusher.Flush()
	}
	return nil
}
