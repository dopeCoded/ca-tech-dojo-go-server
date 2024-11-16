package service

import (
	"math/rand"
	"time"

	"my-go-project/internal/repository"
)

// GachaService はガチャ関連のビジネスロジックを定義するインターフェースです。
type GachaService interface {
	DrawGacha(userID int64, times int) ([]GachaResult, error)
	ListCharacters(userID int64) ([]UserCharacterResponse, error)
}

// GachaResult はガチャの結果を表す構造体です。
type GachaResult struct {
	CharacterID int64  `json:"characterID"`
	Name        string `json:"name"`
}

// UserCharacterResponse はユーザーが所持するキャラクター情報を表す構造体です。
type UserCharacterResponse struct {
	UserCharacterID int64  `json:"userCharacterID"`
	CharacterID     int64  `json:"characterID"`
	Name            string `json:"name"`
}

// gachaService は GachaService インターフェースを実装する構造体です。
type gachaService struct {
	repo repository.GachaRepository
}

// NewGachaService は新しい GachaService を生成します。
func NewGachaService(repo repository.GachaRepository) GachaService {
	return &gachaService{repo}
}

// DrawGacha は指定された回数だけガチャを引き、その結果を返します。
func (s *gachaService) DrawGacha(userID int64, times int) ([]GachaResult, error) {
	// ガチャアイテムと総確率を取得
	items, totalProbability, err := s.repo.GetGachaItems()
	if err != nil {
		return nil, err
	}

	// シードされた乱数ジェネレーターを使用
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	var results []GachaResult
	var characterIDs []int64
	for i := 0; i < times; i++ {
		r := rnd.Float64() * totalProbability
		var cumulative float64
		for _, item := range items {
			cumulative += item.Probability
			if r <= cumulative {
				name, err := s.repo.GetCharacterName(item.CharacterID)
				if err != nil {
					name = "Unknown"
				}
				results = append(results, GachaResult{
					CharacterID: item.CharacterID,
					Name:        name,
				})
				characterIDs = append(characterIDs, item.CharacterID)
				break
			}
		}
	}

	// ユーザーにキャラクターを追加
	if len(characterIDs) > 0 {
		if err := s.repo.AddUserCharacters(userID, characterIDs); err != nil {
			return nil, err
		}
	}

	return results, nil
}

// ListCharacters は指定されたユーザーが所持するキャラクターの一覧を取得します。
func (s *gachaService) ListCharacters(userID int64) ([]UserCharacterResponse, error) {
	// ユーザーが所持するキャラクターを取得
	userCharacters, err := s.repo.GetUserCharacters(userID)
	if err != nil {
		return nil, err
	}

	var responses []UserCharacterResponse
	for _, uc := range userCharacters {
		name, err := s.repo.GetCharacterName(uc.CharacterID)
		if err != nil {
			name = "Unknown"
		}
		responses = append(responses, UserCharacterResponse{
			UserCharacterID: uc.ID,
			CharacterID:     uc.CharacterID,
			Name:            name,
		})
	}

	return responses, nil
}
