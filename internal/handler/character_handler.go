package handler

import (
	"encoding/json"
	"net/http"

	"my-go-project/internal/service"
	"my-go-project/pkg/middleware"
)

// ListCharacters はユーザーが所持するキャラクター一覧を取得します。
func (h *GachaHandler) ListCharacters(w http.ResponseWriter, r *http.Request) {
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

	characters, err := h.gachaService.ListCharacters(userID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	res := struct {
		Characters []service.UserCharacterResponse `json:"characters"`
	}{
		Characters: characters,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
