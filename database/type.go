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
