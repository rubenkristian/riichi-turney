package database

import (
	"fmt"
	"time"
)

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
