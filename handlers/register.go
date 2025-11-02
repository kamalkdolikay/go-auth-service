package handlers

import (
	"auth/db"
	"auth/models"
	"encoding/json"
	"net/http"
	"unicode"

	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// === Request & Response ===
type registerRequest struct {
	Name     string `json:"name" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type registerResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type fieldError struct {
	Field   string
	Message string
}

func (e fieldError) Error() string {
	return e.Message
}

// RegisterHandler handles POST /register
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Parse JSON body
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Invalid JSON", "", http.StatusBadRequest)
		return
	}

	// 2. Normalize email
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	// 3. Validate
	if err := validateRegister(req); err != nil {
		// Handle structured field errors
		if fe, ok := err.(fieldError); ok {
			jsonError(w, fe.Message, strings.ToLower(fe.Field), http.StatusBadRequest)
			return
		}

		// Handle validator.ValidationErrors → collect ALL
		if fieldErrs, ok := err.(validator.ValidationErrors); ok {
			errors := make([]errorResponse, 0, len(fieldErrs))

			for _, e := range fieldErrs {
				msg := fieldMessage(e.Field(), e.Tag(), e.Param())
				errors = append(errors, errorResponse{
					Error: msg,
					Field: strings.ToLower(e.Field()),
				})
			}

			// Send array of errors
			jsonErrors(w, errors, http.StatusBadRequest)
			return
		}

		// Fallback
		jsonError(w, err.Error(), "", http.StatusBadRequest)
		return
	}

	// 4. Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		jsonError(w, "Failed to process password", "", http.StatusInternalServerError)
		return
	}

	// 5. Insert into DB
	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashed),
	}

	id, err := insertUser(user)
	if err != nil {
		if isPostgresUniqueViolation(err) {
			jsonError(w, "Email already taken", "email", http.StatusConflict)
		} else {
			jsonError(w, "Failed to create user", "", http.StatusInternalServerError)
		}
		return
	}

	// 6. Success response
	resp := registerResponse{
		ID:    id,
		Name:  user.Name,
		Email: user.Email,
	}
	jsonResponse(w, resp, http.StatusCreated)
}

// === Validation ===
func validateRegister(r registerRequest) error {
	// Tag-based validation
	if err := validate.Struct(r); err != nil {
		return err
	}

	// Custom: uppercase + digit
	if !hasUpperCase(r.Password) || !hasDigit(r.Password) {
		return fieldError{
			Field:   "Password",
			Message: "password must contain at least one uppercase letter and one digit",
		}
	}
	return nil
}

func hasUpperCase(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

func hasDigit(s string) bool {
	for _, r := range s {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

// === DB ===
func insertUser(u models.User) (int, error) {
	var id int
	query := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id`
	err := db.GetDB().QueryRow(query, u.Name, u.Email, u.Password).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Detect unique violation (email already exists)
func isPostgresUniqueViolation(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505" // unique_violation
	}
	return strings.Contains(err.Error(), "duplicate key")
}
