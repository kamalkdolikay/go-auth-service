// handlers/profile.go
package handlers

import (
	"net/http"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := GetUserFromContext(r)
	if !ok {
		jsonError(w, "User not found in context", "", http.StatusUnauthorized)
		return
	}

	jsonResponse(w, map[string]any{
		"id":    claims.UserID,
		"email": claims.Email,
		"message": "Welcome to your profile!",
	}, http.StatusOK)
}