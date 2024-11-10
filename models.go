package main

import "time"

type User struct {
    ID        int64
    Name      string
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Character struct {
    ID        int64
    Name      string
    Rarity    int
    CreatedAt time.Time
    UpdatedAt time.Time
}

type GachaProbability struct {
    ID           int64
    CharacterID  int64
    Probability  float64
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

type UserCharacter struct {
    ID          int64
    UserID      int64
    CharacterID int64
    AcquiredAt  time.Time
}
