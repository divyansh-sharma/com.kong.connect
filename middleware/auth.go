package middleware

import (
	"context"
	"net/http"
	"strings"
)

// UserContextKey is used to store user info in request context
type contextKey string

const UserContextKey = contextKey("user")

type UserClaims struct {
	Username string
	Roles    []string
}

// Dummy token validation — replace with real JWT validation
func validateToken(token string) (*UserClaims, error) {
	// This is where you'd parse and validate a JWT or token
	if token == "admin-token" {
		return &UserClaims{Username: "admin", Roles: []string{"admin"}}, nil
	}
	if token == "viewer-token" {
		return &UserClaims{Username: "viewer", Roles: []string{"viewer"}}, nil
	}
	return nil, http.ErrNoCookie
}

// AuthMiddleware authenticates requests and injects user info into context
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		user, err := validateToken(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RoleAuthorization checks if user has required role(s)
func RoleAuthorization(allowedRoles ...string) func(http.Handler) http.Handler {
	roleSet := make(map[string]struct{})
	for _, role := range allowedRoles {
		roleSet[role] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(UserContextKey).(*UserClaims)
			if !ok || user == nil {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			for _, role := range user.Roles {
				if _, ok := roleSet[role]; ok {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "Forbidden", http.StatusForbidden)
		})
	}
}
