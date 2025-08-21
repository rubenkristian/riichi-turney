package database

import (
	"os"
	"path/filepath"

	"github.com/rubenkristian/riichi-turney/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DatabaseGame struct {
	db *gorm.DB
}

func CreateDatabaseGame() (*DatabaseGame, error) {
	cfgDir, _ := utils.GetConfigDir("riichi-bot")
	os.MkdirAll(cfgDir, 0755)

	db, err := gorm.Open(sqlite.Open(filepath.Join(cfgDir, "./turney.db")), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	tables := []any{
		&Tournament{},
		&RegisterTournament{},
		&Player{},
		&Match{},
		&PlayerMatch{},
		&TournamentMatch{},
		&TournamentMatchPlayer{},
	}

	if err := db.AutoMigrate(tables...); err != nil {
		return nil, err
	}

	return &DatabaseGame{
		db: db,
	}, nil
}
