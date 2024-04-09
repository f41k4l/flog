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
	loki := make(chan [][]string, 1)
	if config.Loki.URL != "" {
		if config.Loki.Client == nil {
			config.Loki.Client = http.DefaultClient
		}
		go func(config LokiConfig, queue chan [][]string, out io.Writer) {
			for msg := range queue {
				err := config.writeToLoki(msg)
				if err != nil {
					<-queue
					close(queue)
					_, _ = out.Write([]byte("[LOKI LOG FAILED] " + err.Error() + "\n"))
					break
				}
			}
		}(config.Loki, loki, output)
	} else {
		close(loki)
	}

	// Initialize Teams
	teams := make(chan string, 1)
	if config.Teams.Webhook != "" {
		if config.Teams.Client == nil {
			config.Teams.Client = http.DefaultClient
		}
		go func(config TeamsConfig, queue chan string, out io.Writer) {
			for msg := range queue {
				err := config.writeToTeams(msg)
				if err != nil {
					<-queue
					close(queue)
					_, _ = out.Write([]byte("[TEAMS LOG FAILED] " + err.Error() + "\n"))
					break
				}
			}
		}(config.Teams, teams, output)
	} else {
		close(teams)
	}

	return &Writer{
		output: output,
		loki:   loki,
		teams:  teams,
	}, nil
}
