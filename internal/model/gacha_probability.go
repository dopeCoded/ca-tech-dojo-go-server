package model

import "time"

// GachaProbability represents the probability of obtaining a specific character in the gacha.
type GachaProbability struct {
    ID           int64     `json:"id"`
    CharacterID  int64     `json:"character_id"`
    Probability  float64   `json:"probability"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}
