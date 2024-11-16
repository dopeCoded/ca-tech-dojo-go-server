package repository

import (
	"database/sql"
	"time"

	"my-go-project/internal/model"
)

// GachaRepository はガチャ関連のデータベース操作を定義するインターフェースです。
type GachaRepository interface {
	GetGachaItems() ([]model.GachaProbability, float64, error)
	AddUserCharacters(userID int64, characterIDs []int64) error
	GetUserCharacters(userID int64) ([]model.UserCharacter, error)
	GetCharacterName(characterID int64) (string, error)
}

// gachaRepository は GachaRepository インターフェースを実装する構造体です。
type gachaRepository struct {
	db *sql.DB
}

// GetGachaItems はガチャに使用されるキャラクターとその確率を取得します。
// 戻り値にはキャラクターのリストと全確率の合計が含まれます。
func (r *gachaRepository) GetGachaItems() ([]model.GachaProbability, float64, error) {
	rows, err := r.db.Query(`
		SELECT character_id, probability
		FROM gacha_probabilities
	`)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []model.GachaProbability
	var totalProbability float64
	for rows.Next() {
		var item model.GachaProbability
		if err := rows.Scan(&item.CharacterID, &item.Probability); err != nil {
			return nil, 0, err
		}
		totalProbability += item.Probability
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return items, totalProbability, nil
}

// AddUserCharacters はユーザーが取得したキャラクターを user_characters テーブルに追加します。
func (r *gachaRepository) AddUserCharacters(userID int64, characterIDs []int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO user_characters (user_id, character_id, acquired_at)
		VALUES (?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	currentTime := time.Now()
	for _, cid := range characterIDs {
		_, err := stmt.Exec(userID, cid, currentTime)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetCharacterName は指定されたキャラクターIDに対応するキャラクター名を取得します。
func (r *gachaRepository) GetCharacterName(characterID int64) (string, error) {
	var name string
	err := r.db.QueryRow(`
		SELECT name
		FROM characters
		WHERE id = ?
	`, characterID).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}
