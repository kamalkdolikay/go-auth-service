package handlers

import (
	"fmt"
	"io"
	"net/http"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "You reached /get")
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusInternalServerError)
		return
	}

	// Print the received body on the server console
	// fmt.Println("Received:", string(body))
	fmt.Fprintf(w, "Data received: %s", body)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "404 - Page Not Found: %s", r.URL.Path)
}