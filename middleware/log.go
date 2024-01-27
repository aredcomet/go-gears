package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	LevelPanic = "panic"
	LevelFatal = "fatal"
	LevelError = "error"
	LevelWarn  = "warn"
	LevelInfo  = "info"
	LevelTrace = "trace"
)

type LoggerConfig struct {
	Logger *logrus.Logger
}

var DefaultLoggerConfig *LoggerConfig

func init() {
	DefaultLoggerConfig = NewLoggerConfig(LevelInfo)
}

func NewLoggerConfig(logLevel string) *LoggerConfig {
	logLevelMap := map[string]logrus.Level{
		LevelPanic: logrus.PanicLevel,
		LevelFatal: logrus.FatalLevel,
		LevelError: logrus.ErrorLevel,
		LevelWarn:  logrus.WarnLevel,
		LevelInfo:  logrus.InfoLevel,
		LevelTrace: logrus.TraceLevel,
	}

	level, ok := logLevelMap[logLevel]
	if !ok {
		level = logrus.InfoLevel
	}

	logger := CreateLogger(level)
	return &LoggerConfig{Logger: logger}
}

func CreateLogger(logLevel logrus.Level) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logLevel)
	return logger
}

func LoggerMiddleware(next http.Handler, lc *LoggerConfig) http.Handler {
	if lc == nil {
		lc = DefaultLoggerConfig
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &LoggingResponseWriter{w, http.StatusOK, 0, start}
		next.ServeHTTP(lrw, r)

		fields := logrus.Fields{
			"method":     r.Method,
			"requestURI": r.RequestURI,
			"remoteAddr": r.RemoteAddr,
			"status":     lrw.statusCode,
			"latency":    time.Since(lrw.StartedAt).String(),
		}

		lc.Logger.WithFields(fields).Info("Handled request")
	})
}

type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	length     int
	StartedAt  time.Time
}

func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *LoggingResponseWriter) Write(b []byte) (int, error) {
	length, err := lrw.ResponseWriter.Write(b)
	lrw.length += length
	return length, err
}
