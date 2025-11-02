package handlers

import (
    "context"
    "net/http"
    "strings"
)

type contextKey string
const UserContextKey contextKey = "user"

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        auth := r.Header.Get("Authorization")
        if !strings.HasPrefix(auth, "Bearer ") {
            jsonError(w, "Missing or invalid Authorization header", "", http.StatusUnauthorized)
            return
        }

        tokenStr := strings.TrimPrefix(auth, "Bearer ")
        claims, err := ParseJWT(tokenStr)
        if err != nil {
            jsonError(w, "Invalid token", "", http.StatusUnauthorized)
            return
        }

        ctx := context.WithValue(r.Context(), UserContextKey, claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func GetUserFromContext(r *http.Request) (*Claims, bool) {
    claims, ok := r.Context().Value(UserContextKey).(*Claims)
    return claims, ok
}