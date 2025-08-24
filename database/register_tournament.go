package database

import (
	"fmt"
	"strings"
	"time"
)

func (dg *DatabaseGame) SyncPointTournamentPlayerRegister(tournamentId uint64) error {
	sql := `
        UPDATE register_tournaments rt
        SET point = COALESCE((
            SELECT SUM(tmp.final_point)
            FROM tournament_match_players tmp
            JOIN tournament_matches tm ON tm.id = tmp.tournament_match_id
            WHERE tmp.player_id = rt.player_id
              AND tm.tournament_id = rt.tournament_id
        ), 0)
    `
	return dg.db.Exec(sql).Error
}

func (dg *DatabaseGame) ListRegisterTournamentPlayers(isMatch bool, tournamentId uint64, search string, pagination Pagination) ([]RegisterTournament, error) {
	var registeredTournament []RegisterTournament
	query := dg.db.Preload("Player").Preload("Tournament").Model(&RegisterTournament{}).Where("tournament_id = ?", tournamentId)

	if search != "" {
		query = query.Joins("JOIN players ON players.id = register_tournaments.player_id").
			Where("players.riichi_city_name LIKE ?", "%"+search+"%")
	}

	if isMatch {
		query = query.Where(`
		    NOT EXISTS (
		        SELECT 1
		        FROM player_matches pm
		        JOIN matches m ON m.id = pm.match_id
		        WHERE pm.player_id = register_tournaments.player_id
		          AND m.tournament_id = register_tournaments.tournament_id
		          AND m.status != 1
		    )
		`)
	}

	// map user input -> real DB column
	allowedSortBy := map[string]string{
		"id":         "register_tournaments.id",
		"created_at": "register_tournaments.created_at",
		"point":      "register_tournaments.point",
	}

	sortBy, ok := allowedSortBy[pagination.SortBy]
	if !ok {
		sortBy = "register_tournaments.id"
	}

	if strings.ToUpper(pagination.Sort) != "DESC" {
		pagination.Sort = "ASC"
	}

	page := pagination.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pagination.Size

	query = query.Order(fmt.Sprintf("%s %s", sortBy, pagination.Sort)).
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
