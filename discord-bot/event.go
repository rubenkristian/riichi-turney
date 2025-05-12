package discordbot

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// handle when command /register selected
func (db *DiscordBot) EventRegister(event *events.ApplicationCommandInteractionCreate) {
	event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("✅").Build())
}

// handle when command /start-table selected
func (db *DiscordBot) EventStartTable(event *events.ApplicationCommandInteractionCreate) {

}

// handle when command /check-table selected
func (db *DiscordBot) EventCheckTable(event *events.ApplicationCommandInteractionCreate) {

}

// handle when command /check-schedule selected
func (db *DiscordBot) EventCheckSchedule(event *events.ApplicationCommandInteractionCreate) {

}

// handle when command /check-point selected
func (db *DiscordBot) EventCheckPoint(event *events.ApplicationCommandInteractionCreate) {

}
