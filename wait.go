package flog

import "time"

func (w *Writer) Wait() {
	for len(w.teams) > 0 {
		time.Sleep(time.Millisecond)
	}
}
