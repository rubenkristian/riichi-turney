package database

import "fmt"

func (dg *DatabaseGame) GetPlayerByRiichiId(id uint64) (*Player, error) {
	var player Player

	if err := dg.db.First(&player, id).Error; err != nil {
		return nil, err
	}

	return &player, nil
}

func (dg *DatabaseGame) GetPlayerByDiscordId(id uint64) (*Player, error) {
	var player Player

	if err := dg.db.Where("discord_id = ?", id).First(&player).Error; err != nil {
		return nil, err
	}

	return &player, nil
}

func (dg *DatabaseGame) CreatePlayer(body PlayerBody) (*Player, error) {
	playerExists, err := dg.GetPlayerByDiscordId(body.DiscordId)

	if err == nil && playerExists != nil {
		return nil, fmt.Errorf("Failed to register, this discord account already register")
	}

	player := &Player{
		Id:             body.RiichiCityId,
		DiscordName:    body.DiscordName,
		RiichiCityName: body.RiichiCityName,
		DiscordId:      body.DiscordId,
	}

	if err := dg.db.Create(player).Error; err != nil {
		return nil, err
	}

	return player, nil
}

func (dg *DatabaseGame) ListPlayer(player PaginationPlayer) ([]Player, error) {
	var players []Player
	query := dg.db.Model(&players)

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
		fmt.Sprintf("%s %s", player.Pagination.SortBy, player.Pagination.Sort),
	).
		Limit(player.Pagination.Size).
		Offset((max(player.Pagination.Page, 1) - 1) * player.Pagination.Size)

	if err := query.Find(&players).Error; err != nil {
		return nil, err
	}

	return players, nil
}

func (dg *DatabaseGame) DeletePlayer(playerId int) (bool, error) {
	err := dg.db.Delete(&Player{}, playerId).Error

	if err != nil {
		return false, err
	}

	return true, nil
}
