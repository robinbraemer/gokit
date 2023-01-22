# gokit

Development Kit contains utilities useful for production Go applications.

**Provided packages:**

- `log` - logging package using the standardized [logr.Logger](https://github.com/go-logr/logr) interface and https://github.com/uber-go/zap as a backend.
- `otel` - OpenTelemetry instrumentation utilities suggesting [HonneyComb](https://www.honeycomb.io/) as a backend.
- `configutil` - utilities for loading configuration files and environment variables using [Viper](https://github.com/spf13/viper) as backend.
- `httputil` - HTTP utilities like error handler, graceful server shutdown, etc.
- `sync/rungroup` - run goroutines gracefully by waiting for them to finish.
