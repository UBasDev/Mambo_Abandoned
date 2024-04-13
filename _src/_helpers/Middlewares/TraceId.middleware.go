package middlewares

import (
	"context"
	"net/http"
	"time"

	enums "github.com/UBasDev/Mambo/_src/MamboCoreService/Core/MamboCoreService.Application/Enums"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type TraceIdKey struct{}

type ZapLogger struct{}

func TraceIdMiddleware(next http.Handler, logger *zap.Logger, environment enums.Environment) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		generatedTraceId := uuid.New().String()
		ctx := context.WithValue(r.Context(), TraceIdKey{}, generatedTraceId)
		ctx = context.WithValue(ctx, ZapLogger{}, logger.With(zap.String("traceid", generatedTraceId), zap.String("environment", environment.String()), zap.Int64("time", time.Now().Unix()), zap.String("method", r.Method), zap.String("path", r.URL.Path)))
		r = r.WithContext(ctx)
		w.Header().Set("Content-Type", "application/json")
		// Call the next handler, which can be another middleware in the chain or the final handler
		next.ServeHTTP(w, r)
	})
}
