package database

import "fmt"

func (dg *DatabaseGame) CreateMatch(body MatchBody) (*Match, error) {
	match := &Match{
		TournamentId: body.TournamentId,
		Status:       0,
	}

	err := dg.db.Create(match).Error

	if err != nil {
		return nil, err
	}

	return match, nil
}

func (dg *DatabaseGame) GetMatchById(matchId uint64) (*Match, error) {
	var match Match

	if err := dg.db.Preload("Tournament").Preload("PlayerMatches.Player").First(&match, matchId).Error; err != nil {
		return nil, err
	}

	return &match, nil
}

func (dg *DatabaseGame) DeleteMatch(matchId uint64) (bool, error) {
	err := dg.db.Delete(&Match{}, matchId).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (dg *DatabaseGame) UpdateStatusMatch(matchId uint64, status int) error {
	result := dg.db.Model(&Match{}).Where("id = ?", matchId).Update("status", status)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("No rows updated")
	} else {
		return nil
	}
}

func (dg *DatabaseGame) ListNotStartedMatch() ([]Match, error) {
	var matches []Match
	query := dg.db.Preload("PlayerMatches.Player").Model(&Match{}).Where("status != ?", 1)

	if err := query.Find(&matches).Error; err != nil {
		return nil, err
	}

	return matches, nil
}

func (dg *DatabaseGame) ListMatch(match PaginationMatch) ([]Match, error) {
	var matches []Match
	query := dg.db.Model(&Match{})

	if match.Pagination.Search != "" {
		searchQuery := fmt.Sprintf("%%%s%%", match.Pagination.Search)
		query = query.Where("match_name LIKE ?", searchQuery)
	}

	if match.Pagination.SortBy == "" {
		match.Pagination.SortBy = "id"
	}

	if match.Pagination.Sort == "" {
		match.Pagination.Sort = "ASC"
	}

	query = query.Order(
		fmt.Sprintf("%s %s", match.Pagination.SortBy, match.Pagination.Sort),
	).
		Limit(match.Pagination.Size).
		Offset((max(match.Pagination.Page, 1) - 1) * match.Pagination.Size)

	if err := query.Find(&matches).Error; err != nil {
		return nil, err
	}

	return matches, nil
}

func (dg *DatabaseGame) ListMatchByPlayerId(playerId int, match PaginationMatch) ([]Match, error) {
	var matchs []Match
	query := dg.db.Model(&Match{}).Where("player_id = ?", playerId)

	if match.Pagination.Search != "" {
		searchQuery := fmt.Sprintf("%%%s%%", match.Pagination.Search)
		query = query.Where("match_name LIKE ?", searchQuery)
	}

	if match.Pagination.SortBy == "" {
		match.Pagination.SortBy = "id"
	}

	if match.Pagination.Sort == "" {
		match.Pagination.Sort = "ASC"
	}

	query = query.Order(
		fmt.Sprintf("%s %s", match.Pagination.SortBy, match.Pagination.Sort),
	).
		Limit(match.Pagination.Size).
		Offset((max(match.Pagination.Page, 1) - 1) * match.Pagination.Size)

	if err := query.Find(&matchs).Error; err != nil {
		return nil, err
	}

	return matchs, nil
}
