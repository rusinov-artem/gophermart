package middleware

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type ResponseSpy struct {
	w          http.ResponseWriter
	body       bytes.Buffer
	statusCode int
}

func (r *ResponseSpy) Header() http.Header {
	return r.w.Header()
}

func (r *ResponseSpy) Write(bytes []byte) (int, error) {
	return r.body.Write(bytes)
}

func (r *ResponseSpy) WriteHeader(statusCode int) {
	r.statusCode = statusCode
}

var Logger = func(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqBody, _ := io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(reqBody))
			start := time.Now()
			spy := &ResponseSpy{w: w}
			h.ServeHTTP(spy, r)
			w.WriteHeader(spy.statusCode)
			_, _ = w.Write(spy.body.Bytes())
			dur := time.Since(start)
			logger.Info("handling request",
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
				zap.String("req.body", string(reqBody)),
				zap.String("resp.body", spy.body.String()),
				zap.Any("resp.headers", spy.Header()),
				zap.Duration("duration", dur),
				zap.Int("statusCode", spy.statusCode),
				zap.Int("size", spy.body.Len()),
			)
		})
	}
}
