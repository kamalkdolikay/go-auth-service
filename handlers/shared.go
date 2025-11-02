// handlers/common.go
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

// === Shared Validator ===
var validate = validator.New()

// === Shared Response Types ===
type errorResponse struct {
	Error string `json:"error"`
	Field string `json:"field,omitempty"`
}

// === HTTP Helpers ===
func jsonError(w http.ResponseWriter, msg, field string, code int) {
	resp := errorResponse{Error: msg}
	if field != "" {
		resp.Field = field
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(resp)
}

func jsonErrors(w http.ResponseWriter, errs []errorResponse, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(errs)
}

func jsonResponse(w http.ResponseWriter, data any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(data)
}

// === Field Error Messages ===
func fieldMessage(field, tag, param string) string {
	switch field {
	case "Name":
		switch tag {
		case "required":
			return "name is required"
		case "min":
			return fmt.Sprintf("name must be at least %s characters", param)
		}
	case "Email":
		switch tag {
		case "required":
			return "email is required"
		case "email":
			return "invalid email format"
		}
	case "Password":
		switch tag {
		case "required":
			return "password is required"
		case "min":
			return fmt.Sprintf("password must be at least %s characters", param)
		}
	}
	return fmt.Sprintf("invalid %s", strings.ToLower(field))
}