package discordbot

import "github.com/disgoorg/disgo/events"

func (db *DiscordBot) CommandAssignRole(event *events.MessageCreate) {
	// event.Client().Rest().AddMemberRole(*event.GuildID, )
}

func (db *DiscordBot) CommandStartCheckTournament(event *events.MessageCreate) {

}

func (db *DiscordBot) CommandEndCheckTournament(event *events.MessageCreate) {

}

func (db *DiscordBot) CommandCheckMatchNow(event *events.MessageCreate) {

}
