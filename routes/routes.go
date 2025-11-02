package routes

import (
	"auth/handlers"
	"net/http"
	"github.com/gorilla/mux"
)

// func RegisterRoutes() http.Handler {
// 	mux := http.NewServeMux()

// 	// Public
// 	mux.HandleFunc("GET /$", handlers.HelloHandler)
// 	mux.HandleFunc("GET /get", handlers.GetHandler)
// 	mux.HandleFunc("POST /post", handlers.PostHandler)
// 	mux.HandleFunc("POST /register", handlers.RegisterHandler)
// 	mux.HandleFunc("POST /login", handlers.LoginHandler)
// 	mux.HandleFunc("POST /logout", handlers.LogoutHandler)

// 	// Protected routes
// 	protected := http.NewServeMux()
// 	protected.HandleFunc("GET /profile", handlers.ProfileHandler)

// 	// Apply auth middleware
// 	mux.Handle("/api/", http.StripPrefix("/api", handlers.AuthMiddleware(protected)))

// 	return mux
// }

// RegisterRoutesToMux registers *all* routes on the supplied mux.
// This is called from api/index.go.
func RegisterRoutesToMux(r *mux.Router) {
	// Public
	r.HandleFunc("/", handlers.HelloHandler).Methods("GET")
	r.HandleFunc("/get", handlers.GetHandler).Methods("GET")
	r.HandleFunc("/post", handlers.PostHandler).Methods("POST")
	r.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	r.HandleFunc("/logout", handlers.LogoutHandler).Methods("POST")
	r.HandleFunc("/health", handlers.HealthHandler).Methods("GET")

	// Protected
	protected := mux.NewRouter()
	protected.HandleFunc("/profile", handlers.ProfileHandler).Methods("GET")
	authHandler := handlers.AuthMiddleware(protected)
	r.PathPrefix("/api/").Handler(http.StripPrefix("/api", authHandler))
}
