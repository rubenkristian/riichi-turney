package database

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func (dg *DatabaseGame) GetTournamentMatchById(id string) (*TournamentMatch, error) {
	var tournamentMatch TournamentMatch

	if err := dg.db.Preload("Tournament").Preload("TournamentMatchPlayers").First(&tournamentMatch, id).Error; err != nil {
		return nil, err
	}

	return &tournamentMatch, nil
}

func (dg *DatabaseGame) CreateTournamentMatch(body TournamentMatchBody) (*TournamentMatch, error) {
	tx := dg.db.Begin()

	if tx.Error != nil {
		return nil, tx.Error
	}

	var existing TournamentMatch
	err := dg.db.First(&existing, body.RoomID).Error

	if err == nil {
		return nil, fmt.Errorf("an active tournament already exists")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	tournamentMatch := TournamentMatch{
		Id:           body.RoomID,
		PaiPuId:      body.PaiPuId,
		TournamentId: body.TournamentId,
	}

	if err := tx.Create(&tournamentMatch).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, player := range body.Players {
		tournamentMatchPlayer := TournamentMatchPlayer{
			PlayerId:          player.PlayerId,
			Score:             player.Score,
			TournamentMatchId: tournamentMatch.Id,
		}

		if err := tx.Create(&tournamentMatchPlayer).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	var loadedTournamentMatch TournamentMatch
	if err := dg.db.Preload("Tournament").
		Preload("TournamentMatchPlayers").
		First(&loadedTournamentMatch, "id = ?", tournamentMatch.Id).Error; err != nil {
		return nil, err
	}

	return &loadedTournamentMatch, nil
}

func (dg *DatabaseGame) ListTournamentMatch(tournament PaginationTournament) ([]TournamentMatch, error) {
	var tournamentMatches []TournamentMatch
	query := dg.db.Model(&Tournament{})

	if tournament.FromDate != nil && tournament.ToDate != nil {
		query = query.Where("created_at BETWEEN ? AND ?", tournament.FromDate, tournament.ToDate)
	}

	// Whitelist sort fields
	allowedSortBy := map[string]bool{"id": true, "created_at": true}
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

	if err := query.Find(&tournamentMatches).Error; err != nil {
		return nil, err
	}

	return tournamentMatches, nil
}
