package model

import "time"

// Character represents a character that can be obtained via gacha.
type Character struct {
    ID        int64     `json:"id"`
    Name      string    `json:"name"`
    Rarity    int       `json:"rarity"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
