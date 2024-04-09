package flog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func (w *Writer) Write(p []byte) (n int, err error) {

	if w.loki != nil {
		err := w.writeToLoki(p)
		if err != nil {
			w.output.Write([]byte("[LOKI FAILED] " + err.Error() + "\n"))
		}
	}

	if w.teams != nil {
		if len(w.teams) > 0 {
			data := <-w.teams
			data = append(data, []byte(fmt.Sprintf("<br /><code>%s</code>", p))...)
			w.teams <- data
		} else {
			w.teams <- []byte(fmt.Sprintf("<code>%s</code>", p))
		}
	}

	n, err = w.output.Write(p)
	if err != nil {
		return
	}

	return
}

func (w *Writer) writeToLoki(p []byte) (err error) {

	buffer := new(bytes.Buffer)
	err = json.NewEncoder(buffer).Encode(lokiWriter{
		Streams: []stream{{
			Stream: w.loki.Labels,
			Values: [][]string{{fmt.Sprint(time.Now().UnixNano()), string(p)}},
		}},
	})
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, w.loki.URL+w.loki.Path, buffer)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	if w.loki.Headers != nil {
		for k, v := range w.loki.Headers {
			if len(v) == 0 {
				continue
			}
			req.Header.Set(k, v[0])
		}
	}

	resp, err := w.loki.Client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if status := resp.StatusCode; status < 200 || status >= 300 {
		d, _ := io.ReadAll(resp.Body)
		err = fmt.Errorf("unexpected status code %d\n%s", status, d)
	}

	return
}

func (config *TeamsConfig) writeToTeams(p []byte) (err error) {

	buffer := new(bytes.Buffer)
	err = json.NewEncoder(buffer).Encode(teamsMessage{
		Title: config.Title,
		Text:  string(p),
	})
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, config.Webhook, buffer)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := config.Client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if status := resp.StatusCode; status < 200 || status >= 300 {
		d, _ := io.ReadAll(resp.Body)
		err = fmt.Errorf("unexpected status code %d\n%s", status, d)
	}

	return
}
