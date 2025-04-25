package database

type Player struct {
	Id             int64
	DiscordName    string
	RiichiCityName string
	DiscordId      int64
	RiichiCityId   int64
}

type Table struct {
	Id        int64
	TableName string
	Day       int
}

type PlayerTable struct {
	TableId  int64
	PlayerId int64
	Score    int
}

type Team struct {
	Id   int64
	Name string
}

type PlayerTeam struct {
	PlayerId int64
	Position string
	IsLeader bool
}
