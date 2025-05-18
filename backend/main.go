package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/golang-jwt/jwt/v4"
)

var secretKey = []byte("secretkey123")

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var (
	items []Item
	idSeq int
	mu    sync.Mutex
)

func main() {
	http.HandleFunc("/api/login", withCORS(loginHandler))
	http.HandleFunc("/api/items", withCORS(authMiddleware(itemsHandler)))
	// Untuk detail item (update dan delete)
	http.HandleFunc("/api/items/", withCORS(authMiddleware(itemDetailHandler)))

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Middleware CORS sederhana
func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Boleh semua origin untuk development, sesuaikan di produksi
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// Simple login: cek email & password tetap, return JWT token
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	fmt.Printf("Login attempt: %s / %s\n", creds.Email, creds.Password)

	// Dummy valid user
	if creds.Email == "user@example.com" && creds.Password == "password" {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"email": creds.Email,
		})
		tokenString, err := token.SignedString(secretKey)
		if err != nil {
			http.Error(w, "Could not create token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
		return
	}

	http.Error(w, "Invalid credentials", http.StatusUnauthorized)
}

// Middleware cek token JWT di header Authorization Bearer
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}
		var tokenString string
		fmt.Sscanf(authHeader, "Bearer %s", &tokenString)

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method")
			}
			return secretKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// Handler CRUD items (GET all & POST create)
func itemsHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(items)
	case http.MethodPost:
		var newItem Item
		err := json.NewDecoder(r.Body).Decode(&newItem)
		if err != nil || strings.TrimSpace(newItem.Name) == "" {
			http.Error(w, "Bad request: missing or invalid name", http.StatusBadRequest)
			return
		}
		idSeq++
		newItem.ID = idSeq
		items = append(items, newItem)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newItem)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Handler update dan delete item by ID: /api/items/{id}
func itemDetailHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	// Ambil ID dari URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[3])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPut:
		var updated Item
		err := json.NewDecoder(r.Body).Decode(&updated)
		if err != nil || strings.TrimSpace(updated.Name) == "" {
			http.Error(w, "Bad request: missing or invalid name", http.StatusBadRequest)
			return
		}
		for i, item := range items {
			if item.ID == id {
				items[i].Name = updated.Name
				json.NewEncoder(w).Encode(items[i])
				return
			}
		}
		http.Error(w, "Item not found", http.StatusNotFound)

	case http.MethodDelete:
		for i, item := range items {
			if item.ID == id {
				items = append(items[:i], items[i+1:]...)
				w.WriteHeader(http.StatusNoContent) // 204 No Content
				return
			}
		}
		http.Error(w, "Item not found", http.StatusNotFound)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
