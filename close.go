package flog

import (
	"time"
)

func (w *Writer) Close() error {
	for len(w.teams) > 0 || len(w.loki) > 0 {
		time.Sleep(50 * time.Millisecond)
	}
	return nil
}
