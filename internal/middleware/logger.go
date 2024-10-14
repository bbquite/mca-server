package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type ResponseWriter interface {
	Header() http.Header
	Write([]byte) (int, error)
	WriteHeader(statusCode int)
}

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	//r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func RequestsLoggingMiddleware(logger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			responseData := &responseData{
				status: 0,
				size:   0,
			}
			loggingWriter := loggingResponseWriter{
				ResponseWriter: w,
				responseData:   responseData,
			}

			h.ServeHTTP(&loggingWriter, r)

			duration := time.Since(start)

			logger.Infoln(
				"method", r.Method,
				"status", responseData.status,
				"uri", r.RequestURI,
				"size", responseData.size,
				"duration", duration,
				//"headers", r.Header,
			)
		})
	}
}
