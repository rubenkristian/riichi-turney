package discordbot

import (
	"fmt"
	"strconv"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// handle when command /register selected
func (db *DiscordBot) EventRegister(event *events.ApplicationCommandInteractionCreate) {
	data := event.SlashCommandInteractionData()

	riichi_city_id := data.String("riichi_city_id")

	riichiId, err := strconv.ParseUint(riichi_city_id, 10, 64)

	if err != nil {
		event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Please to insert riichi id as number").Build())
	}

	playerGame, err := db.DbGame.GetPlayerByRiichiId(riichiId)

	if err != nil {
		fmt.Println(err.Error())
	}

	if playerGame != nil {
		event.CreateMessage(discord.NewMessageCreateBuilder().SetContent(fmt.Sprintf("This id %s already registered", riichi_city_id)).Build())
	}

	player, err := db.RiichiCommand.FindPlayer(riichi_city_id)

	if err != nil {
		fmt.Println(err.Error())
		event.CreateMessage(discord.NewMessageCreateBuilder().SetContent(fmt.Sprintf("This id %s not found", riichi_city_id)).Build())
		return
	}

	if len(player.FriendList) <= 0 {
		event.CreateMessage(discord.NewMessageCreateBuilder().SetContent(fmt.Sprintf("This id %s not found", riichi_city_id)).Build())
		return
	}

	playerData := player.FriendList[0]

	confirmBtn := discord.NewPrimaryButton("Confirm", fmt.Sprintf("confirm:%s,%d", playerData.Nickname, playerData.UserID))
	cancelBtn := discord.NewSecondaryButton("Cancel", "cancel")

	event.CreateMessage(
		discord.NewMessageCreateBuilder().
			SetContent(fmt.Sprintf("Username: %s, Riichi City Id: %s ?", playerData.Nickname, riichi_city_id)).
			AddActionRow(confirmBtn, cancelBtn).Build(),
	)
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
