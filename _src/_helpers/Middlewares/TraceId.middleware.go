package middlewares

import (
	"context"
	"net/http"

	"github.com/openzipkin/zipkin-go"
	"go.uber.org/zap"
)

type TraceIdKey struct{}

type ZapLogger struct{}

func TraceIdMiddleware(next http.Handler, logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span := zipkin.SpanFromContext(r.Context())
		//span.Tag("traceid", span.Context().ID.String()) Mevcut spane key-value şeklinde tag eklememizi sağlar.
		//span.Annotate(time.Now(), "annotation1")        Mevcut spane key-value şeklinde annotation eklememizi sağlar.
		generatedTraceId := span.Context().TraceID.String()
		ctx := context.WithValue(r.Context(), TraceIdKey{}, generatedTraceId)
		ctx = context.WithValue(ctx, ZapLogger{}, logger.With(zap.String("traceid", generatedTraceId), zap.String("method", r.Method), zap.String("path", r.URL.Path)))
		r = r.WithContext(ctx)
		w.Header().Set("Content-Type", "application/json")
		// Call the next handler, which can be another middleware in the chain or the final handler
		next.ServeHTTP(w, r)
	})
}
