package database

import "time"

type Tournament struct {
	Id          uint64     `gorm:"primaryKey" json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	StartAt     time.Time  `json:"start_at"`
	EndAt       time.Time  `json:"end_at"`
	Active      bool       `json:"active"`
	RoleID      string     `json:"role_id"`
	RegisterEnd *time.Time `json:"register_end"`
	ClassifyID  string     `json:"classify_id"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"index"`

	// Added relationships
	Registrations []RegisterTournament `gorm:"foreignKey:tournament_id"`
	Matches       []Match              `gorm:"foreignKey:tournament_id"`
}

type RegisterTournament struct {
	Id           uint64     `gorm:"primaryKey" json:"id"`
	PlayerId     uint64     `json:"player_id"`
	TournamentId uint64     `json:"tournament_id"`
	Point        int64      `json:"point"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    *time.Time `json:"deleted_at" gorm:"index"`

	// Added relationships
	Player     Player     `gorm:"foreignKey:player_id"`
	Tournament Tournament `gorm:"foreignKey:tournament_id"`
}

type Player struct {
	Id             uint64     `gorm:"primaryKey" json:"id"`
	DiscordName    string     `json:"discord_name"`
	RiichiCityName string     `json:"riichi_city_name"`
	DiscordId      uint64     `json:"discord_id"`
	CreatedAt      time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt      *time.Time `json:"deleted_at" gorm:"index"`

	// Added relationships
	Registrations []RegisterTournament `gorm:"foreignKey:player_id"`
	PlayerMatches []PlayerMatch        `gorm:"foreignKey:player_id"`
}

type Match struct {
	Id           uint64     `gorm:"primaryKey" json:"id"`
	TournamentId uint64     `json:"tournament_id"`
	Status       int        `json:"status"` // 0 created, -1 failed to start, 1 playing
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    *time.Time `json:"deleted_at" gorm:"index"`

	// Added relationships
	Tournament    Tournament    `gorm:"foreignKey:tournament_id"`
	PlayerMatches []PlayerMatch `gorm:"foreignKey:match_id"`
}

type PlayerMatch struct {
	Id       uint64 `gorm:"primaryKey" json:"id"`
	MatchId  uint64 `json:"match_id"`
	PlayerId uint64 `json:"player_id"`

	// Added relationships
	Player Player `gorm:"foreignKey:player_id"`
}

type TournamentMatch struct {
	Id           string     `gorm:"primaryKey" json:"id"`
	PaiPuId      string     `json:"pai_pu_id"`
	TournamentId uint64     `json:"tournament_id"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    *time.Time `json:"deleted_at" gorm:"index"`

	Tournament             Tournament              `gorm:"foreignKey:tournament_id"`
	TournamentMatchPlayers []TournamentMatchPlayer `gorm:"foreignKey:tournament_match_id"`
}

type TournamentMatchPlayer struct {
	Id                uint64 `gorm:"primaryKey" json:"id"`
	TournamentMatchId string `json:"tournament_match_id"`
	Score             int64  `json:"score"`
	Point             int64  `json:"point"`
	Penalty           int64  `json:"penalty"`
	FinalPoint        int64  `json:"final_point"`
	PlayerId          uint64 `json:"player_id"`
	Status            bool   `json:"status"`

	Player Player `gorm:"foreignKey:player_id"`
}
