package handler

import (
	"encoding/json"
	"net/http"

	"my-go-project/internal/service"
	"my-go-project/pkg/auth"
	"my-go-project/pkg/middleware"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService}
}

// CreateUser は新しいユーザーを作成し、JWT トークンを返します。
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Bad Request: name is required", http.StatusBadRequest)
		return
	}

	user, err := h.userService.CreateUser(req.Name)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	res := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// GetUser はユーザー情報を取得します。
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// コンテキストからユーザーIDを取得
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.userService.GetUser(userID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	res := struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}{
		ID:   user.ID,
		Name: user.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// UpdateUser はユーザー名を更新します。
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// コンテキストからユーザーIDを取得
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Bad Request: name is required", http.StatusBadRequest)
		return
	}

	if err := h.userService.UpdateUser(userID, req.Name); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User updated successfully"))
}
