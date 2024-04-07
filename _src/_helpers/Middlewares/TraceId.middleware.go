package middlewares

import (
	"context"
	"net/http"

	models "github.com/UBasDev/Mambo/_src/MamboCoreService/Core/MamboCoreService.Application/Models"
	"github.com/google/uuid"
)

func TraceIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), models.TraceIdKey{}, uuid.New().String())

		r = r.WithContext(ctx)
		w.Header().Set("Content-Type", "application/json")
		// Call the next handler, which can be another middleware in the chain or the final handler
		next.ServeHTTP(w, r)
	})
}
