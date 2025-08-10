package discordbot

import (
	"fmt"
	"strconv"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/rubenkristian/riichi-turney/database"
)

// handle when command /register selected
func (db *DiscordBot) EventRegister(event *events.ApplicationCommandInteractionCreate) {
	user := event.User()
	discordID := uint64(user.ID)
	discordName := user.Username

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

	if playerGame != nil {
		confirmBtn := discord.NewPrimaryButton("Confirm", fmt.Sprintf("confirm:%d", playerGame.Id))
		cancelBtn := discord.NewSecondaryButton("Cancel", "cancel")

		event.CreateMessage(
			discord.NewMessageCreateBuilder().
				SetContent(fmt.Sprintf("Username: %s, Riichi City Id: %s ?", playerGame.RiichiCityName, riichi_city_id)).
				AddActionRow(confirmBtn, cancelBtn).Build(),
		)
		return
	}

	riichiCityId, err := strconv.ParseUint(riichi_city_id, 10, 64)
	if err != nil {
		event.CreateMessage(
			discord.NewMessageCreateBuilder().SetContent("❌ Registration cancelled, riichi city id not number").
				SetEphemeral(true).
				Build(),
		)
	}

	newPlayer, err := db.DbGame.CreatePlayer(database.PlayerBody{
		DiscordId:      discordID,
		DiscordName:    discordName,
		RiichiCityName: playerData.Nickname,
		RiichiCityId:   riichiCityId,
	})

	if err != nil {
		event.CreateMessage(
			discord.NewMessageCreateBuilder().SetContent(fmt.Sprintf("❌ Registration cancelled, %s", err.Error())).
				SetEphemeral(true).
				Build(),
		)
	}

	confirmBtn := discord.NewPrimaryButton("Confirm", fmt.Sprintf("confirm:%d", newPlayer.Id))
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
