package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DatabaseGame struct {
	db *gorm.DB
}

func CreateDatabaseGame() (*DatabaseGame, error) {
	db, err := gorm.Open(sqlite.Open("./turney.db"), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	tables := []any{
		&Tournament{},
		&RegisterTournament{},
		&Player{},
		&Match{},
		&PlayerMatch{},
		&Point{},
	}

	if err := db.AutoMigrate(tables...); err != nil {
		return nil, err
	}

	return &DatabaseGame{
		db: db,
	}, nil
}
