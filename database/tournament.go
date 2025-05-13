package database

import "fmt"

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
	var tournament *Tournament

	err := dg.db.Where("active = ?", true).First(tournament).Error

	if err != nil {
		return nil, err
	}

	tournament = &Tournament{
		Name:        body.Name,
		Description: body.Description,
		StartAt:     body.StartAt,
		EndAt:       body.EndAt,
		RegisterEnd: body.RegisterEnd,
	}

	err = dg.db.Create(tournament).Error

	if err != nil {
		return nil, err
	}

	return tournament, nil
}

func (dg *DatabaseGame) ListTournament(tournament PaginationTournament) ([]Tournament, error) {
	var tournaments []Tournament
	query := dg.db.Model(&Tournament{})

	if tournament.Pagination.Search != "" {
		query = query.Where("name LIKE ?", fmt.Sprintf("%%%s%%", tournament.Pagination.Search))
	}

	if tournament.FromDate != nil && tournament.ToDate != nil {
		query = query.Where("craeted_at BETWEEN ?,?", tournament.FromDate, tournament.ToDate)
	}

	if tournament.Pagination.SortBy == "" {
		tournament.Pagination.SortBy = "id"
	}

	if tournament.Pagination.Sort == "" {
		tournament.Pagination.Sort = "ASC"
	}

	query = query.Order(
		fmt.Sprintf("%s %s", tournament.Pagination.SortBy, tournament.Pagination.Sort),
	).
		Limit(tournament.Pagination.Size).
		Offset((max(tournament.Pagination.Page, 1) - 1) * tournament.Pagination.Size)

	if err := query.Find(&tournaments).Error; err != nil {
		return nil, err
	}

	return tournaments, nil
}
