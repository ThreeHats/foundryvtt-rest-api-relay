package middleware

import (
	"fmt"
	"net/http"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
)

// RequireScope returns middleware that checks the request context for a required scope.
// Master key users (ScopedKey == nil) always pass through.
// Scoped keys must have the required scope in their scope list.
func RequireScope(scope string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqCtx := helpers.GetRequestContext(r)
			if reqCtx == nil {
				helpers.WriteError(w, http.StatusUnauthorized, "Authentication required")
				return
			}

			// Master key users have all scopes
			if reqCtx.ScopedKey == nil {
				next.ServeHTTP(w, r)
				return
			}

			if !reqCtx.ScopedKey.HasScope(scope) {
				helpers.WriteError(w, http.StatusForbidden,
					fmt.Sprintf("API key lacks required scope: %s", scope))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
