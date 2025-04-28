package discordbot

import "github.com/disgoorg/disgo/events"

func (db *DiscordBot) CommandAssignRole(event *events.MessageCreate) {
	// event.Client().Rest().AddMemberRole(*event.GuildID, )
}

func (db *DiscordBot) CommandCheckTable(event *events.MessageCreate) {

}

func (db *DiscordBot) CommandCheckSelf(event *events.MessageCreate) {

}

func (db *DiscordBot) CommandCheckRank(event *events.MessageCreate) {

}

func (db *DiscordBot) CommandCheckOwnTable(event *events.MessageCreate) {

}

// for admin role
func (db *DiscordBot) CommandStartTable(event *events.MessageCreate) {

}
