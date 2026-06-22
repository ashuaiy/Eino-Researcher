package utils

import (
	"encoding/json"
	"fmt"
	"io"
)

type SSEWriter struct {
	w io.Writer
}

func NewSSEWriter(w io.Writer) SSEWriter {
	return SSEWriter{w: w}
}

func (s SSEWriter) Event(event string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(s.w, "event: %s\ndata: %s\n\n", event, data)
	return err
}
