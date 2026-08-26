package middleware

import (
	"context"
	"net/http"
	"strings"

	"product-inventory/internal/auth"
	sharedHttp "product-inventory/internal/http"
)

type contextKey string

const (
	userClaimsKey contextKey = "user_claims"
)

// Authenticate is a middleware that parses and validates a Bearer JWT token from the Authorization header.
func Authenticate(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				sharedHttp.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authorization header is missing")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				sharedHttp.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authorization header must be Bearer token")
				return
			}

			tokenStr := parts[1]
			claims, err := auth.ValidateToken(tokenStr, secret)
			if err != nil {
				sharedHttp.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid or expired authorization token")
				return
			}

			// Store claims in request context
			ctx := context.WithValue(r.Context(), userClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserClaims extracts the authenticated user claims from the context.
func GetUserClaims(ctx context.Context) (*auth.UserClaims, bool) {
	claims, ok := ctx.Value(userClaimsKey).(*auth.UserClaims)
	return claims, ok
}

// RequireRole is a middleware that checks if the authenticated user has one of the required roles.
// This middleware MUST run after Authenticate.
func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := GetUserClaims(r.Context())
			if !ok {
				sharedHttp.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
				return
			}

			roleAllowed := false
			for _, role := range allowedRoles {
				if strings.EqualFold(claims.Role, role) {
					roleAllowed = true
					break
				}
			}

			if !roleAllowed {
				sharedHttp.WriteError(w, http.StatusForbidden, "FORBIDDEN_ACCESS", "You do not have permission to access this resource")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
