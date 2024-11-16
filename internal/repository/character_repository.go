package repository

import (
	"database/sql"

	"my-go-project/internal/model"
)

// NewGachaRepository は新しい GachaRepository を生成します。
func NewGachaRepository(db *sql.DB) GachaRepository {
	return &gachaRepository{db}
}

// GetUserCharacters は指定されたユーザーが所持するキャラクターを取得します。
func (r *gachaRepository) GetUserCharacters(userID int64) ([]model.UserCharacter, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, character_id, acquired_at
		FROM user_characters
		WHERE user_id = ?
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userCharacters []model.UserCharacter
	for rows.Next() {
		var uc model.UserCharacter
		if err := rows.Scan(&uc.ID, &uc.UserID, &uc.CharacterID, &uc.AcquiredAt); err != nil {
			return nil, err
		}
		userCharacters = append(userCharacters, uc)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userCharacters, nil
}