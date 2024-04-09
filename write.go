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

	// Write to loki
	select {
	case data, ok := <-w.loki:
		if ok {
			data = append(data, []string{fmt.Sprint(time.Now().UnixNano()), string(p)})
			w.loki <- data
		}
	default:
		w.loki <- [][]string{{fmt.Sprint(time.Now().UnixNano()), string(p)}}
	}

	// Write to teams
	select {
	case data, ok := <-w.teams:
		if ok {
			data += fmt.Sprintf("<pre>%s</pre>", p)
			w.teams <- data
		}
	default:
		w.teams <- fmt.Sprintf("<pre>%s</pre>", p)
	}

	n, err = w.output.Write(p)
	if err != nil {
		return
	}

	return
}

func (config *LokiConfig) writeToLoki(p [][]string) (err error) {

	buffer := new(bytes.Buffer)
	err = json.NewEncoder(buffer).Encode(lokiWriter{
		Streams: []stream{{
			Stream: config.Labels,
			Values: p,
		}},
	})
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, config.URL+config.Path, buffer)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	if config.Headers != nil {
		for k, v := range config.Headers {
			if len(v) == 0 {
				continue
			}
			req.Header.Set(k, v[0])
		}
	}

	resp, err := config.Client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if status := resp.StatusCode; status < 200 || status >= 300 {
		d, _ := io.ReadAll(resp.Body)
		err = fmt.Errorf("unexpected status code %d: %s", status, d)
	}

	return
}

func (config *TeamsConfig) writeToTeams(p string) (err error) {

	buffer := new(bytes.Buffer)
	err = json.NewEncoder(buffer).Encode(teamsMessage{
		Title: config.Title,
		Text:  p,
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
		err = fmt.Errorf("unexpected status code %d: %s", status, d)
	}

	return
}
