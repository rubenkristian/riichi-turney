package database

import "time"

type Tournament struct {
	Id          uint       `gorm:"primaryKey" json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	StartAt     time.Time  `json:"start_at"`
	EndAt       time.Time  `json:"end_at"`
	Active      bool       `json:"active"`
	RegisterEnd *time.Time `json:"register_end"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   time.Time  `json:"deleted_at" gorm:"index"`
}

type RegisterTournament struct {
	Id           uint       `gorm:"primaryKey" json:"id"`
	PlayerId     uint       `json:"player_id"`
	Player       Player     `gorm:"foreignKey:PlayerId" json:"player"`
	TournamentId uint       `json:"tournament_id"`
	Tournament   Tournament `json:"tournament"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    time.Time  `json:"deleted_at" gorm:"index"`
}

type Player struct {
	Id             uint      `gorm:"primaryKey" json:"id"`
	DiscordName    string    `json:"discord_name"`
	RiichiCityName string    `json:"riichi_city_name"`
	DiscordId      int64     `json:"discord_id"`
	RiichiCityId   int64     `json:"riichi_city_id"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt      time.Time `json:"deleted_at" gorm:"index"`
}

type Match struct {
	Id           uint       `gorm:"primaryKey" json:"id"`
	TableName    string     `json:"table_name"`
	Day          time.Time  `json:"day"`
	TournamentId uint       `json:"tournament_id"`
	Tournament   Tournament `json:"tournament"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    time.Time  `json:"deleted_at" gorm:"index"`
}

type PlayerMatch struct {
	Id       uint    `gorm:"primaryKey" json:"id"`
	MatchId  uint    `json:"match_id"`
	Match    Match   `gorm:"foreignKey:MatchId" json:"match"`
	PlayerId uint    `json:"player_id"`
	Player   Player  `gorm:"foreignKey:PlayerId" json:"player"`
	Score    float64 `json:"score"`
}

type Point struct {
	Id           uint       `gorm:"primaryKey" json:"id"`
	Value        float64    `json:"value"`
	Type         bool       `json:"type"`
	PlayerId     uint       `json:"player_id"`
	Player       Player     `gorm:"foreignKey:PlayerId" json:"player"`
	TournamentId uint       `json:"tournament_id"`
	Tournament   Tournament `json:"tournament"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    time.Time  `json:"deleted_at" gorm:"index"`
}
