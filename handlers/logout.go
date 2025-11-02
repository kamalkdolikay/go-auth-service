// handlers/logout.go
package handlers

import "net/http"

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
    jsonResponse(w, map[string]string{"message": "Logged out"}, http.StatusOK)
}