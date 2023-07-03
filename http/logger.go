package http

import (
	"log"
	"net/http"
	"time"
)

func Logs(logger *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lrw := loggingResponseWriter{
			ResponseWriter: w, // compose original http.ResponseWriter
			status:         0,
			size:           0,
		}

		next.ServeHTTP(&lrw, r)

		dur := time.Since(start)

		log.Printf("Request %s %s ~ status=%d took=%v size=%d", r.Method, r.RequestURI, lrw.status, dur, lrw.size)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter // compose original http.ResponseWriter
	status              int
	size                int
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b) // write response using original http.ResponseWriter
	r.size += size                         // capture size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode) // write status code using original http.ResponseWriter
	r.status = statusCode                    // capture status code
}
