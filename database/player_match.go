package database

func (dg *DatabaseGame) CreatePlayerMatch(matchId uint64, playerId uint64) (PlayerMatch, error) {

	return PlayerMatch{}, nil
}

func (dg *DatabaseGame) UpdatePlayerMatch() (PlayerMatch, error) {
	return PlayerMatch{}, nil
}

func (dg *DatabaseGame) GetPlayerInMatch(matchId int) ([]PlayerMatch, error) {
	return []PlayerMatch{}, nil
}
