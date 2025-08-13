package database

import "fmt"

func (dg *DatabaseGame) PointsByPlayer(playerId uint64) ([]TournamentMatchPlayer, error) {
	var points []TournamentMatchPlayer
	query := dg.db.Model(&TournamentMatchPlayer{}).Where("player_id != ?", playerId)

	if err := query.Find(&points).Error; err != nil {
		return nil, err
	}

	return points, nil
}

func (dg *DatabaseGame) UpdatePointTournamentMatchPlayer(tournamentMatchId string, playerId uint64, point int64) error {
	result := dg.db.
		Model(&TournamentMatchPlayer{}).
		Where("tournament_match_id = ? AND player_id = ?", tournamentMatchId, playerId).Update("point", point)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("No rows updated")
	} else {
		return nil
	}
}

func (dg *DatabaseGame) UpdateFinalPointTournamentMatchPlayer(tournamentMatchPlayerId uint64, penalty int64) error {
	var tournamentMatchPlayer TournamentMatchPlayer

	if err := dg.db.First(&tournamentMatchPlayer, tournamentMatchPlayerId).Error; err != nil {
		return err
	}

	finalPoint := tournamentMatchPlayer.Point - penalty

	result := dg.db.
		Model(&TournamentMatchPlayer{}).
		Where("id = ?", tournamentMatchPlayerId).Updates(map[string]any{
		"penalty":     penalty,
		"final_point": finalPoint,
	})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("No rows updated")
	} else {
		return nil
	}
}
