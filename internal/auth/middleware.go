package auth

import (
	"context"
	"net/http"

	clerk "github.com/clerk/clerk-sdk-go/v2"
	clerkHttp "github.com/clerk/clerk-sdk-go/v2/http"
)

// A private key for context that only this package can access.
type contextKey string

const UserContextKey = contextKey("user")

// Middleware decodes the session token and adds the user ID to the context.
func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		protected := clerkHttp.RequireHeaderAuthorization()
		return protected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := clerk.SessionClaimsFromContext(r.Context())
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			// Token is valid, add claims to context
			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}))
	}
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *clerk.SessionClaims {
	raw, _ := ctx.Value(UserContextKey).(*clerk.SessionClaims)
	return raw
}
