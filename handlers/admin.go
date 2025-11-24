package handlers

import (
	"auth/db"
	"net/http"
	"time"
)

// AdminDashboardResponse defines the JSON structure for the frontend
type AdminDashboardResponse struct {
	TotalUsers    int           `json:"total_users"`
	TotalRequests int           `json:"total_requests"`
	Users         []UserStat    `json:"users"`
	RecentLogs    []ActivityLog `json:"recent_logs"`
}

type UserStat struct {
	ID           int        `json:"id"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	Role         string     `json:"role"`
	RequestCount int        `json:"request_count"`
	LastActive   *time.Time `json:"last_active"` // Pointer allows null
}

type ActivityLog struct {
	ID        int       `json:"id"`
	UserEmail string    `json:"user_email"`
	Prompt    string    `json:"prompt"`
	Topic     string    `json:"topic"`
	Type      string    `json:"type"`
	Time      time.Time `json:"time"`
}

// AdminDashboardHandler returns metrics for the Admin Dashboard
func AdminDashboardHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Check if requester is Admin
	claims, ok := GetUserFromContext(r)
	if !ok {
		jsonError(w, "Unauthorized", "", http.StatusUnauthorized)
		return
	}

	// Check role in DB to be sure (JWT might be old)
	var role string
	err := db.GetDB().QueryRow("SELECT role FROM users WHERE id = $1", claims.UserID).Scan(&role)
	if err != nil || role != "admin" {
		jsonError(w, "Access Denied: Admins Only", "", http.StatusForbidden)
		return
	}

	// 2. Get User Stats (Aggregation)
	// Lists all users, counts their activities, and finds their last active timestamp
	rows, err := db.GetDB().Query(`
		SELECT 
			u.id, u.name, u.email, u.role,
			COUNT(a.id) as request_count,
			MAX(a.created_at) as last_active
		FROM users u
		LEFT JOIN user_activity a ON u.id = a.user_id
		GROUP BY u.id
		ORDER BY last_active DESC NULLS LAST
	`)
	if err != nil {
		jsonError(w, "Database error fetching users", "", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []UserStat
	totalReqs := 0

	for rows.Next() {
		var u UserStat
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.RequestCount, &u.LastActive); err != nil {
			continue
		}
		users = append(users, u)
		totalReqs += u.RequestCount
	}

	// 3. Get Recent Global Activity Logs
	logRows, err := db.GetDB().Query(`
		SELECT a.id, u.email, a.prompt, a.topic_detected, a.request_type, a.created_at
		FROM user_activity a
		JOIN users u ON a.user_id = u.id
		ORDER BY a.created_at DESC
		LIMIT 20
	`)
	if err != nil {
		jsonError(w, "Database error fetching logs", "", http.StatusInternalServerError)
		return
	}
	defer logRows.Close()

	var logs []ActivityLog
	for logRows.Next() {
		var l ActivityLog
		// Handle potential NULLs if data is messy, but schema defines most as NOT NULL usually
		if err := logRows.Scan(&l.ID, &l.UserEmail, &l.Prompt, &l.Topic, &l.Type, &l.Time); err != nil {
			continue
		}
		logs = append(logs, l)
	}

	// 4. Return Data
	resp := AdminDashboardResponse{
		TotalUsers:    len(users),
		TotalRequests: totalReqs,
		Users:         users,
		RecentLogs:    logs,
	}

	jsonResponse(w, resp, http.StatusOK)
}
