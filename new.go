package flog

import (
	"io"
	"net/http"
	"os"
)

// New creates a new flog.Writer with the given configuration.
//
// Exxample usage:
//
//	    l, err := flog.New(flog.Config{
//	        Output: os.Stdout,
//	        Loki: flog.LokiConfig{
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

	// Initialize Loki
	loki := &config.Loki
	if loki.URL == "" {
		loki = nil
	} else {
		if loki.Client == nil {
			loki.Client = http.DefaultClient
		}
	}

	// Initialize Teams
	teamsConfig := &config.Teams
	var teams chan []byte
	if teamsConfig.Webhook == "" {
		teamsConfig = nil
	} else {
		teams = make(chan []byte, 1)
		if teamsConfig.Client == nil {
			teamsConfig.Client = http.DefaultClient
		}
		go func(config *TeamsConfig, queue chan []byte, out io.Writer) {
			for msg := range queue {
				err := config.writeToTeams(msg)
				if err != nil {
					_, _ = out.Write([]byte("[TEAMS FAILED] " + err.Error()))
					close(queue)
					return
				}
			}
		}(teamsConfig, teams, output)
	}

	return &Writer{
		output: output,
		loki:   loki,
		teams:  teams,
	}, nil
}
