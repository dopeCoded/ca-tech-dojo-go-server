package handler

import (
	"encoding/json"
	"net/http"

	"my-go-project/internal/service"
	"my-go-project/pkg/middleware"
)

type GachaHandler struct {
	gachaService service.GachaService
}

func NewGachaHandler(gachaService service.GachaService) *GachaHandler {
	return &GachaHandler{gachaService}
}

// DrawGacha はガチャを引くリクエストを処理します。
func (h *GachaHandler) DrawGacha(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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
		Times int `json:"times"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Times <= 0 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	results, err := h.gachaService.DrawGacha(userID, req.Times)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	res := struct {
		Results []service.GachaResult `json:"results"`
	}{
		Results: results,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

