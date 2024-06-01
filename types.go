package flog

import (
	"io"
	"net/http"
)

type Writer struct {
	defaultOutput io.Writer
	loki          chan [][]string
	teams         chan string
}

type Config struct {
	Output io.Writer
	Loki   LokiConfig
	Teams  TeamsConfig
}

type LokiConfig struct {
	URL       string
	Path      string
	Headers   http.Header
	Labels    map[string]string
	BasicAuth BasicAuth
	Client    *http.Client
}

type TeamsConfig struct {
	Title   string
	Webhook string
	Client  *http.Client
}

type BasicAuth struct {
	Username string
	Password string
}

type lokiWriter struct {
	Streams []lokiStream `json:"streams"`
}

type lokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

type teamsMessage struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}
