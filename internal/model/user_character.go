package model

import "time"

// UserCharacter represents a character that a user has obtained.
type UserCharacter struct {
    ID          int64     `json:"user_character_id"`
    UserID      int64     `json:"user_id"`
    CharacterID int64     `json:"character_id"`
    AcquiredAt  time.Time `json:"acquired_at"`
}
