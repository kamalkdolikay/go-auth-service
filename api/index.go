// api/index.go
package handler

import (
	"net/http"
	"os"

	"auth/db"
	"auth/handlers"
	"auth/routes"

	"github.com/gorilla/mux"
)

// Vercel calls this function on every request
func Handler(w http.ResponseWriter, r *http.Request) {
	// Initialise once per cold-start
	handlers.InitJWT()
	db.InitDB()

	// Build the router (same as in routes.go)
	router := mux.NewRouter()
	routes.RegisterRoutesToMux(router)

	// Serve
	router.ServeHTTP(w, r)
}

// Vercel requires the handler to be exported as a top-level func
func main() {
	// Vercel will never call main(); we keep it for local testing
	http.HandleFunc("/", Handler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, nil)
}
