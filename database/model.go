package database

import "time"

type Tournament struct {
	Id          uint64     `gorm:"primaryKey" json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	StartAt     time.Time  `json:"start_at"`
	EndAt       time.Time  `json:"end_at"`
	Active      bool       `json:"active"`
	RegisterEnd *time.Time `json:"register_end"`
	RoleID      string     `json:"role_id"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   time.Time  `json:"deleted_at" gorm:"index"`

	// Added relationships
	Registrations []RegisterTournament `gorm:"foreignKey:tournament_id"`
	Matches       []Match              `gorm:"foreignKey:tournament_id"`
	Points        []Point              `gorm:"foreignKey:tournament_id"`
}

type RegisterTournament struct {
	Id           uint64    `gorm:"primaryKey" json:"id"`
	PlayerId     uint64    `json:"player_id"`
	TournamentId uint64    `json:"tournament_id"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    time.Time `json:"deleted_at" gorm:"index"`

	// Added relationships
	Player     Player     `gorm:"foreignKey:player_id"`
	Tournament Tournament `gorm:"foreignKey:tournament_id"`
}

type Player struct {
	Id             uint      `gorm:"primaryKey" json:"id"`
	DiscordName    string    `json:"discord_name"`
	RiichiCityName string    `json:"riichi_city_name"`
	DiscordId      uint64    `json:"discord_id"`
	RiichiCityId   uint64    `json:"riichi_city_id"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt      time.Time `json:"deleted_at" gorm:"index"`

	// Added relationships
	Registrations []RegisterTournament `gorm:"foreignKey:player_id"`
	PlayerMatches []PlayerMatch        `gorm:"foreignKey:player_id"`
	Points        []Point              `gorm:"foreignKey:player_id"`
}

type Match struct {
	Id           uint64    `gorm:"primaryKey" json:"id"`
	TableName    string    `json:"table_name"`
	Day          time.Time `json:"day"`
	TournamentId uint      `json:"tournament_id"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    time.Time `json:"deleted_at" gorm:"index"`

	// Added relationships
	Tournament    Tournament    `gorm:"foreignKey:tournament_id"`
	PlayerMatches []PlayerMatch `gorm:"foreignKey:match_id"`
}

type PlayerMatch struct {
	Id       uint64  `gorm:"primaryKey" json:"id"`
	MatchId  uint64  `json:"match_id"`
	PlayerId uint64  `json:"player_id"`
	Score    float64 `json:"score"`

	// Added relationships
	Player Player `gorm:"foreignKey:player_id"`
	Match  Match  `gorm:"foreignKey:match_id"`
}

type Point struct {
	Id           uint64    `gorm:"primaryKey" json:"id"`
	Value        float64   `json:"value"`
	Type         bool      `json:"type"`
	PlayerId     uint64    `json:"player_id"`
	TournamentId uint64    `json:"tournament_id"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    time.Time `json:"deleted_at" gorm:"index"`

	// Added relationships
	Player     Player     `gorm:"foreignKey:player_id"`
	Tournament Tournament `gorm:"foreignKey:tournament_id"`
}
