package flog

import (
	"io"
	"net/http"
)

type Writer struct {
	output io.Writer
	loki   *LokiConfig
	teams  *TeamsConfig
}

type Config struct {
	Output io.Writer
	Loki   LokiConfig
	Teams  TeamsConfig
}

type LokiConfig struct {
	URL     string
	Path    string
	Headers http.Header
	Labels  map[string]string
	Client  *http.Client
}

type TeamsConfig struct {
	Title   string
	Webhook string
	Client  *http.Client
}

type lokiWriter struct {
	Streams []stream `json:"streams"`
}

type stream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

type teamsMessage struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}
