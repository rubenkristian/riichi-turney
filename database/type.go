package database

import "time"

type TournamentBody struct {
	TournamentId uint64    `json:"tournament_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	StartAt      time.Time `json:"start_at"`
	EndAt        time.Time `json:"end_at"`
	RegisterEnd  time.Time `json:"register_end"`
	ClassifyID   string    `json:"classify_id"`
	RoleID       string    `json:"role_id"`
}

type InputTournamentMatchPlayer struct {
	PlayerId uint64 `json:"player_id"`
	Score    int64  `json:"score"`
}

type TournamentMatchBody struct {
	PaiPuId      string                       `json:"pai_pu_id"`
	TournamentId uint64                       `json:"tournament_id"`
	RoomID       string                       `json:"room_id"`
	Players      []InputTournamentMatchPlayer `json:"players"`
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
	DiscordId      uint64 `json:"discord_id"`
	RiichiCityId   uint64 `json:"riichi_city_id"`
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
	TournamentId uint64    `json:"tournament_id"`
	Players      []uint64  `json:"players"`
}

type MatchBodyUpdate struct {
	MatchName    string    `json:"match_name"`
	Day          time.Time `json:"day"`
	TournamentId uint64    `json:"tournament_id"`
	MatchPlayers []uint64  `json:"match_players"`
	Players      []uint64  `json:"players"`
}

type PaginationMatch struct {
	Pagination Pagination
	FromDate   *string
	ToDate     *string
}

type PlayerCheck struct {
	RiichiId  uint64
	DiscordId uint64
	Status    int
}
