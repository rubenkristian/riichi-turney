package database

func (dg *DatabaseGame) CreateMatch(body MatchBody) (*Match, error) {
	match := &Match{
		TableName: body.MatchName,
		Day: body.Day,
		TournamentId: body.TournamentId,
	}

	err := dg.db.Create(match).Error

	if err != nil {
		return nil, err
	}

	return match, nil
}

func (dg *DatabaseGame) DeleteMatch(matchId uint) (bool, error) {
	err := dg.db.Delete(&Match{}, matchId).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (dg *DatabaseGame) ListMatch(match PaginationMatch) ([]Match, error) {
	var matchs []Match
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
		fmt.Sprintf("%s %s", match.Pagination.SortBy, match.Pagination.Sort)
	).
		Limit(match.Pagination.Size).
		Offset((max(match.Pagination.Page, 1) - 1) * match.Pagination.Size)

	if err := query.Find(&matchs).Error; err != nil {
		return nil, err
	}

	return matchs, nil 
}

func (dg *DatabaseGame) DetailMatch(id int) (*Match, error) {
	var match *Match

	if id <= 0 {
		err := dg.db.Order('created_at desc').First(match)

		if err != nil {
			return nil, err
		}

		return match, nil
	}

	err := dg.db.First(match, id).Error

	if err != nil {
		return nil, err
	}

	return match, nil
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
		fmt.Sprintf("%s %s", match.Pagination.SortBy, match.Pagination.Sort)
	).
		Limit(match.Pagination.Size).
		Offset((max(match.Pagination.Page, 1) - 1) * match.Pagination.Size)

	if err := query.Find(&matchs).Error; err != nil {
		return nil, err
	}

	return matchs, nil
}
