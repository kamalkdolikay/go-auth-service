package handlers

import (
	"auth/db"
	"encoding/json"
	"log"
	"net/http"
)

type ActivityEvent struct {
	UserID         int    `json:"user_id"` // Required
	Prompt         string `json:"prompt"`
	TopicDetected  string `json:"topic_detected"`
	TargetLanguage string `json:"target_language"`
	RequestType    string `json:"request_type"` // "voice" | "text"
}

func ActivityLoggedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var event ActivityEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Basic validation
	if event.UserID <= 0 || event.Prompt == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Fire-and-forget insert — we don't care if it fails slightly
	go func() {
		query := `
			INSERT INTO user_activity (user_id, prompt, topic_detected, target_language, request_type)
			VALUES ($1, $2, $3, $4, $5)
		`
		_, err := db.GetDB().Exec(query, event.UserID, event.Prompt, event.TopicDetected, event.TargetLanguage, event.RequestType)
		if err != nil {
			log.Printf("Failed to log activity for user %d: %v", event.UserID, err)
			// Intentionally silent — logging is best-effort
		}
	}()

	// Immediate 200 — do not block AI service
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"acknowledged": true})
}
