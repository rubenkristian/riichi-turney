package database

import "time"

type TournamentBody struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartAt     time.Time `json:"start_at"`
	EndAt       time.Time `json:"end_at"`
	RegisterEnd time.Time `json:"register_end"`
}

type Pagination struct {
	Page   int
	Size   int
	SortBy string
	Sort   string
	Search string
}

type PaginationTournament struct {
	Pagination Pagination
	FromDate   *string
	ToDate     *string
}

type PlayerBody struct {
	DiscordName    string `json:"discord_name"`
	RiichiCityName string `json:"riichi_city_name"`
	DiscordId      int64  `json:"discord_id"`
	RiichiCityId   int64  `json:"riichi_city_id"`
}

type PaginationPlayer struct {
	Pagination Pagination
	FromDate   *string
	ToDate     *string
}

type PaginationPoint struct {
	Pagination Pagination
	FromDate   *string
	ToDate     *string
	PlayerId   int
}

type MatchBody struct {
	MatchName    string    `json:"match_name"`
	Day          time.Time `json:"day"`
	TournamentId uint      `json:"tournament_id"`
}

type PaginationMatch struct {
	Pagination Pagination
	FromDate   *string
	ToDate     *string
}
