package database

import (
	"fmt"
	"strings"
	"time"
)

func (dg *DatabaseGame) ListRegisterTournamentPlayers(pagination Pagination) ([]RegisterTournament, error) {
	var registeredTournament []RegisterTournament
	query := dg.db.Preload("Player").Preload("Tournament").Model(&RegisterTournament{})

	// Whitelist sort fields
	allowedSortBy := map[string]bool{"id": true, "created_at": true}
	if !allowedSortBy[pagination.SortBy] {
		pagination.SortBy = "id"
	}

	if strings.ToUpper(pagination.Sort) != "DESC" {
		pagination.Sort = "ASC"
	}

	page := pagination.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pagination.Size

	query = query.Order(fmt.Sprintf("%s %s", pagination.SortBy, pagination.Sort)).
		Limit(pagination.Size).
		Offset(offset)

	if err := query.Find(&registeredTournament).Error; err != nil {
		return nil, err
	}

	return registeredTournament, nil
}

func (dg *DatabaseGame) GetRegisterTournamentPlayers(tournamentId uint64) ([]uint64, error) {
	var ids []uint64
	err := dg.db.Model(&RegisterTournament{}).Select("player_id").Where("tournament_id = ?", tournamentId).Pluck("player_id", &ids).Error

	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (dg *DatabaseGame) CreateRegisterTournament(playerId uint64) (*RegisterTournament, error) {
	var tournament Tournament

	err := dg.db.Where("active = ?", true).First(&tournament).Error

	if err != nil {
		return nil, err
	}

	if tournament.RegisterEnd != nil && time.Now().Compare(*tournament.RegisterEnd) == 1 {
		return nil, fmt.Errorf("cannot register to this tournament, tournament registration already end")
	}

	var exist RegisterTournament

	if err := dg.db.Where("tournament_id = ? AND player_id = ?", tournament.Id, playerId).First(&exist).Error; err == nil {
		return nil, fmt.Errorf("this player already registered")
	}

	registerTournament := RegisterTournament{
		PlayerId:     playerId,
		TournamentId: tournament.Id,
	}

	if err := dg.db.Create(&registerTournament).Error; err != nil {
		return nil, err
	}

	if err := dg.db.Preload("Player").Preload("Tournament").First(&registerTournament, registerTournament.Id).Error; err != nil {
		return nil, err
	}

	return &registerTournament, nil
}
