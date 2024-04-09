package flog

import (
	"net/http"
	"os"
)

// New creates a new lokilogger.Writer with the given configuration.
//
// Exxample usage:
//
//	    l, err := lokilogger.New(lokilogger.Config{
//	        Output: os.Stdout,
//	        Loki: logger.LokiConfig{
//	            URL:  "http://localhost:3100",
//	            Path: "/loki/api/v1/push",
//	            Labels: map[string]string{
//	                "app": "myapp",
//	            },
//	        },
//	    })
//	    if err != nil {
//	        panic(err)
//	    }
//	    defer l.Close()
//
//			log.SetOutput(l)
//			log.SetReportTimestamp(false)
func New(config Config) (*Writer, error) {
	output := config.Output
	if output == nil {
		output = os.Stdout
	}

	loki := &config.Loki
	if loki.URL == "" {
		loki = nil
	} else {
		if loki.Client == nil {
			loki.Client = http.DefaultClient
		}
	}

	return &Writer{
		output: output,
		loki:   loki,
	}, nil
}
