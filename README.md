# lokilogger

`lokilogger` is a simple logging client in pure Golang. It is designed to be easy to use and to provide a simple interface for logging messages to the console and to loki instance.

# Usage

```go
logger, err := lokilogger.New(lokilogger.Config{
  Output: os.Stdout,
  Loki: logger.LokiConfig{
    URL:  "http://localhost:3100",
    Path: "/loki/api/v1/push",
    Headers: http.Header{
      "X-Scope-OrgID": []string{"1"},
    },
    Labels: map[string]string{
      "app": "myapp",
    },
  },
})
if err != nil {
  panic(err)
}
defer logger.Close()

log.SetOutput(logger)
log.SetReportTimestamp(false)
```
