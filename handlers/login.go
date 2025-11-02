package handlers

import (
	"auth/db"
	"auth/models"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

// === Request & Response ===
type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type loginResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// LoginHandler handles POST /login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Parse JSON
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Invalid JSON", "", http.StatusBadRequest)
		return
	}

	// 2. Normalize email
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	// 3. Validate input
	if err := validate.Struct(req); err != nil {
		if fieldErrs, ok := err.(validator.ValidationErrors); ok {
			resp := make([]errorResponse, 0, len(fieldErrs))
			for _, e := range fieldErrs {
				msg := fieldMessage(e.Field(), e.Tag(), e.Param())
				resp = append(resp, errorResponse{
					Error: msg,
					Field: strings.ToLower(e.Field()),
				})
			}
			jsonErrors(w, resp, http.StatusBadRequest)
			return
		}
		jsonError(w, err.Error(), "", http.StatusBadRequest)
		return
	}

	// 4. Find user
	user, err := getUserByEmail(req.Email)
	if err != nil {
		// Hide existence of email
		jsonError(w, "Invalid email or password", "", http.StatusUnauthorized)
		return
	}

	// 5. Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		jsonError(w, "Invalid email or password", "", http.StatusUnauthorized)
		return
	}

	// 6. Success – return JWT in body (no cookie)
	token, err := GenerateJWT(user.ID, user.Email)
	if err != nil {
		jsonError(w, "Failed to generate token", "", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, map[string]any{
		"access_token": token,
		"token_type":   "bearer",
		"expires_in":   jwtExpiresMinutes * 60,
		"user": loginResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}, http.StatusOK)
}

// getUserByEmail fetches user by normalized email
func getUserByEmail(email string) (models.User, error) {
	var user models.User
	query := `SELECT id, name, email, password FROM users WHERE email = $1`
	err := db.GetDB().QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	return user, err
}
