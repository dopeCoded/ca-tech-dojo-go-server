package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type UserCreateRequest struct {
	Name string `json:"name"`
}

type UserCreateResponse struct {
	Token string `json:"token"`
}

type UserGetResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type UserUpdateRequest struct {
	Name string `json:"name"`
}

// データベースへのユーザ作成
func createUser(name string) (*User, error) {
	result, err := db.Exec("INSERT INTO users (name) VALUES (?)", name)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:   id,
		Name: name,
	}

	return user, nil
}

// `/user/create` ハンドラ
func userCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req UserCreateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		log.Printf("User name is empty")
		http.Error(w, "Bad Request: name is required", http.StatusBadRequest)
		return
	}

	user, err := createUser(req.Name)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	token, err := GenerateJWT(user.ID)
	if err != nil {
		log.Printf("Failed to generate JWT token: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := UserCreateResponse{Token: token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// `/user/get` ハンドラ
func userGetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// 認証トークンの検証
	tokenString := r.Header.Get("x-token")
	if tokenString == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	claims, err := ValidateJWT(tokenString)
	if err != nil {
		log.Printf("Invalid JWT token: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// ユーザ情報の取得
	var user User
	err = db.QueryRow("SELECT id, name FROM users WHERE id = ?", claims.UserID).Scan(&user.ID, &user.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		log.Printf("Failed to get user: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	res := UserGetResponse{
		ID:   user.ID,
		Name: user.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// `/user/update` ハンドラ
func userUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// 認証トークンの検証
	tokenString := r.Header.Get("x-token")
	if tokenString == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	claims, err := ValidateJWT(tokenString)
	if err != nil {
		log.Printf("Invalid JWT token: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req UserUpdateRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		log.Printf("User name is empty")
		http.Error(w, "Bad Request: name is required", http.StatusBadRequest)
		return
	}

	// ユーザ情報の更新
	_, err = db.Exec("UPDATE users SET name = ?, updated_at = NOW() WHERE id = ?", req.Name, claims.UserID)
	if err != nil {
		log.Printf("Failed to update user: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User updated successfully"))
}
