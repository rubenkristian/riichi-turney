package database

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func (dg *DatabaseGame) GetTournament(id uint) (*Tournament, error) {
	var tournament *Tournament
	if id == 0 {
		dg.db.Order("created_at desc").First(tournament)

		return tournament, nil
	}

	err := dg.db.First(tournament, id).Error

	if err != nil {
		return nil, err
	}

	return tournament, nil
}

func (dg *DatabaseGame) CreateTournament(body TournamentBody) (*Tournament, error) {
	var existing Tournament
	err := dg.db.Where("active = ?", true).First(&existing).Error
	if err == nil {
		return nil, fmt.Errorf("an active tournament already exists")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	tournament := Tournament{
		Name:        body.Name,
		Description: body.Description,
		StartAt:     body.StartAt,
		EndAt:       body.EndAt,
		RegisterEnd: &body.RegisterEnd,
		RoleID:      body.RoleID,
		Active:      true,
	}

	if err = dg.db.Create(&tournament).Error; err != nil {
		return nil, err
	}

	return &tournament, nil
}

func (dg *DatabaseGame) ListTournament(tournament PaginationTournament) ([]Tournament, error) {
	var tournaments []Tournament
	query := dg.db.Model(&Tournament{})

	if tournament.Pagination.Search != "" {
		query = query.Where("name LIKE ?", fmt.Sprintf("%%%s%%", tournament.Pagination.Search))
	}

	if tournament.FromDate != nil && tournament.ToDate != nil {
		query = query.Where("created_at BETWEEN ? AND ?", tournament.FromDate, tournament.ToDate)
	}

	// Whitelist sort fields
	allowedSortBy := map[string]bool{"id": true, "name": true, "created_at": true}
	if !allowedSortBy[tournament.Pagination.SortBy] {
		tournament.Pagination.SortBy = "id"
	}

	if strings.ToUpper(tournament.Pagination.Sort) != "DESC" {
		tournament.Pagination.Sort = "ASC"
	}

	page := tournament.Pagination.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * tournament.Pagination.Size

	query = query.Order(fmt.Sprintf("%s %s", tournament.Pagination.SortBy, tournament.Pagination.Sort)).
		Limit(tournament.Pagination.Size).
		Offset(offset)

	if err := query.Find(&tournaments).Error; err != nil {
		return nil, err
	}

	return tournaments, nil
}
