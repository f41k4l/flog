package flog

func (w *Writer) Close() error {
	close(w.teams)
	return nil
}
