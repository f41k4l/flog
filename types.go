package flog

import (
	"io"
	"net/http"
)

type Writer struct {
	output io.Writer
	loki   *LokiConfig
}

type Config struct {
	Output io.Writer
	Loki   LokiConfig
}

type LokiConfig struct {
	URL     string
	Path    string
	Headers http.Header
	Labels  map[string]string
	Client  *http.Client
}

type lokiWriter struct {
	Streams []stream `json:"streams"`
}

type stream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}
