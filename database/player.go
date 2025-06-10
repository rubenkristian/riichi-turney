package database

import "fmt"

func (dg *DatabaseGame) GetPlayer(id uint) (*Player, error) {
	var player *Player

	if id <= 0 {
		dg.db.Order('created_at desc').First(player)

		return player, nil
	}

	err := dg.db.First(player, id).Error

	if err != nil {
		return nil, err
	}

	return player, nil
}

func (dg *DatabaseGame) CreatePlayer(body PlayerBody) {
	player := &Player{
		DiscordName: body.DiscordName,
		RiichiCityName: body.RiichiCityName,
		DiscordId: body.DiscordId,
		RiichiCityId: body.RiichiCityId,
	}

	err := dg.db.Create(player).Error

	if err != nil {
		return nil, err
	}

	return player, nil
}

func (dg *DatabaseGame) ListPlayer(player PaginationPlayer) {
	var players []Player
	query := dg.db.Model(&player{})

	if player.Pagination.Search != "" {
		searchQuery := fmt.Sprintf("%%%s%%", player.Pagination.Search)
		query = query.Where("discord_name LIKE ? OR riichi_city_name LIKE ?", searchQuery, searchQuery)
	}

	if player.Pagination.SortBy == "" {
		player.Pagination.SortBy = "id"
	}

	if player.Pagination.Sort == "" {
		player.Pagination.Sort = "ASC"
	}

	query = query.Order(
		fmt.Sprintf("%s %s", player.Pagination.SortBy, player.Pagination.Sort)
	).
		Limit(player.Pagination.Size).
		Offset((max(player.Pagination.Page, 1) - 1) * player.Pagination.Size)

	if err := query.Find(&players).Error; err != nil {
		return nil, err
	}

	return players, nil
}

func (dg *DatabaseGame) DeletePlayer(int playerId) (bool, error) {
	err := dg.db.Delete(&Player{}, playerId).Error

	if err != nil {
		return false, err
	}

	return true, nil
}
