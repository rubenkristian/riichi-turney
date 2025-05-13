package database

import (
	"fmt"
	"time"
)

func (dg *DatabaseGame) CreateRegisterTournament() (*RegisterTournament, error) {
	var tournament *Tournament

	err := dg.db.Where("active = ?", true).First(tournament).Error

	if err != nil {
		return nil, err
	}

	if time.Now().Compare(tournament.RegisterEnd) == 1 {
		return nil, fmt.Errorf("cannot register to this tournament, tournament registration already end")
	}

	return nil, nil
}
