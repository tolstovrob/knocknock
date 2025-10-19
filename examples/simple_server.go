package main

/*
 * Пример использования knocknock с in-memory реализацией хранилища и тремя ручками. Комментариев не будет, тут не так
 * сложно :)
 */

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tolstovrob/knocknock"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

var (
	auth *knocknock.Auth
)

func main() {
	store := knocknock.HandleMemoryStore()
	auth = knocknock.HandleAuth(store)

	mux := http.NewServeMux()
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/logout", logoutHandler)
	mux.HandleFunc("/profile", profileHandler)

	handler := auth.Middleware()(mux)

	fmt.Println("Server starting on :8080")
	http.ListenAndServe(":8080", handler)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
	}

	session, err := auth.CreateSession(r.Context(), user)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     auth.AuthOptions.CookieName,
		Value:    session.Token,
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Logged in successfully. Token: %s. Do not tell anyone", session.Token)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session := knocknock.GetSession(r.Context())
	if session != nil {
		auth.DeleteSession(r.Context(), session.Token)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     auth.AuthOptions.CookieName,
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Logged out successfully")
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	session := knocknock.GetSession(r.Context())
	if session == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, ok := session.UserData.(User)
	if !ok {
		http.Error(w, "Invalid session data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
