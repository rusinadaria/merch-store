package middleware

import (
    "net/http"
    "time"
    "log/slog"
)

func LoggerMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        lrw := &LoggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

        next.ServeHTTP(lrw, r)

        logger.Info("Request",
            slog.String("remote_addr", r.RemoteAddr),
            slog.String("time", time.Now().Format("02/Jan/2006:15:04:05 -0700")),
            slog.String("method", r.Method),
            slog.String("path", r.URL.Path),
            slog.String("proto", r.Proto),
            slog.Int("status_code", lrw.statusCode),
            slog.Int("content_length", lrw.contentLength),
            slog.String("referer", r.Referer()),
            slog.String("user_agent", r.UserAgent()),
            slog.Int64("duration", time.Since(start).Milliseconds()),
        )
    })
}

func LoggerMiddlewareWrapper(logger *slog.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return LoggerMiddleware(logger, next)
    }
}

type LoggingResponseWriter struct {
    http.ResponseWriter
    statusCode   int
    contentLength int
}

func (lrw *LoggingResponseWriter) WriteHeader(code int) {
    lrw.statusCode = code
    lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *LoggingResponseWriter) Write(b []byte) (int, error) {
    n, err := lrw.ResponseWriter.Write(b)
    lrw.contentLength += n
    return n, err
}
