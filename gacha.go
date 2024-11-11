package main

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// ガチャ実行リクエスト
type GachaDrawRequest struct {
	Times int `json:"times"`
}

// ガチャ実行レスポンス
type GachaDrawResponse struct {
	Results []GachaResult `json:"results"`
}

// ガチャ結果
type GachaResult struct {
	CharacterID int64  `json:"characterID"`
	Name        string `json:"name"`
}

// ユーザキャラクター一覧レスポンス
type CharacterListResponse struct {
	Characters []UserCharacterResponse `json:"characters"`
}

// ユーザキャラクター情報
type UserCharacterResponse struct {
	UserCharacterID int64  `json:"userCharacterID"`
	CharacterID     int64  `json:"characterID"`
	Name            string `json:"name"`
}

// ガチャ実行ハンドラ
func gachaDrawHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	// リクエストボディの読み取り
	var req GachaDrawRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Times <= 0 {
		log.Printf("Invalid GachaDrawRequest: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// ガチャを実行
	results, err := drawGacha(req.Times)
	if err != nil {
		log.Printf("Failed to draw gacha: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// ユーザの所持キャラクターに追加
	err = addUserCharacters(claims.UserID, results)
	if err != nil {
		log.Printf("Failed to add user characters: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// レスポンスを作成
	res := GachaDrawResponse{Results: results}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// ガチャ実行ロジック
func drawGacha(times int) ([]GachaResult, error) {
	// ガチャ確率を取得
	rows, err := db.Query(`
        SELECT gp.character_id, gp.probability, c.name
        FROM gacha_probabilities gp
        JOIN characters c ON gp.character_id = c.id
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type GachaItem struct {
		CharacterID int64
		Probability float64
		Name        string
	}

	var items []GachaItem
	var totalProbability float64
	for rows.Next() {
		var item GachaItem
		if err := rows.Scan(&item.CharacterID, &item.Probability, &item.Name); err != nil {
			return nil, err
		}
		totalProbability += item.Probability
		items = append(items, item)
	}

	if totalProbability <= 0 {
		return nil, errors.New("total probability is zero or negative")
	}

	// ガチャを引く
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	var results []GachaResult
	for i := 0; i < times; i++ {
		r := rnd.Float64() * totalProbability
		var cumulative float64
		for _, item := range items {
			cumulative += item.Probability
			if r <= cumulative {
				results = append(results, GachaResult{
					CharacterID: item.CharacterID,
					Name:        item.Name,
				})
				break
			}
		}
	}

	return results, nil
}

// ユーザの所持キャラクターに追加
func addUserCharacters(userID int64, results []GachaResult) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO user_characters (user_id, character_id, acquired_at) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, result := range results {
		_, err := stmt.Exec(userID, result.CharacterID, time.Now())
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// `/character/list` ハンドラ
func characterListHandler(w http.ResponseWriter, r *http.Request) {
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

	// ユーザの所持キャラクターを取得
	rows, err := db.Query(`
        SELECT uc.id, uc.character_id, c.name
        FROM user_characters uc
        JOIN characters c ON uc.character_id = c.id
        WHERE uc.user_id = ?
    `, claims.UserID)
	if err != nil {
		log.Printf("Failed to query user characters: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var characters []UserCharacterResponse
	for rows.Next() {
		var uc UserCharacterResponse
		if err := rows.Scan(&uc.UserCharacterID, &uc.CharacterID, &uc.Name); err != nil {
			log.Printf("Failed to scan user character: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		characters = append(characters, uc)
	}

	// レスポンスを作成
	res := CharacterListResponse{Characters: characters}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
