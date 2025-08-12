package discordbot

import (
	"fmt"
	"strconv"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
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
	data := event.SlashCommandInteractionData()

	match_id := data.String("match_id")

	matchId, err := strconv.ParseUint(match_id, 10, 64)

	if err != nil {
		event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Please to insert match id as number").Build())
	}

	match, err := db.DbGame.GetMatchById(matchId)

	if err != nil {
		event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Match Not Found").Build())
	}

	lobbyPlayers, err := db.RiichiCommand.FetchTournamentPlayers(match.TournamentId)

	var playerReady []uint64 = []uint64{}
	var playerNotReady []database.PlayerCheck = []database.PlayerCheck{}

	idIndex := 0

	for _, lobbyPlayer := range lobbyPlayers {
		if lobbyPlayer.UserID == match.PlayerMatches[idIndex].PlayerId {
			if lobbyPlayer.Status == 2 {
				playerReady = append(playerReady, lobbyPlayer.UserID)
			} else {
				playerNotReady = append(playerNotReady, database.PlayerCheck{
					RiichiId:  match.PlayerMatches[idIndex].PlayerId,
					DiscordId: match.PlayerMatches[idIndex].Player.DiscordId,
					Status:    lobbyPlayer.Status,
				})
			}
			idIndex += 1
		}
	}

	for idIndex < len(match.PlayerMatches) {
		playerNotReady = append(playerNotReady, database.PlayerCheck{
			RiichiId:  match.PlayerMatches[idIndex].PlayerId,
			DiscordId: match.PlayerMatches[idIndex].Player.DiscordId,
			Status:    -1,
		})
		idIndex += 1
	}

	if len(playerNotReady) > 0 {
		mentions := ""

		for _, player := range playerNotReady {
			mentions += fmt.Sprintf("<@%d>\n", player.DiscordId)
		}
		content := "The following players, please get ready in the tournament lobby — the match is about to start!\n\n" + mentions
		_, err := db.Client.Rest().CreateMessage(
			snowflake.MustParse(db.Setting.ChannelNotify),
			discord.NewMessageCreateBuilder().SetContent(content).Build(),
		)

		if err != nil {
			event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("some player not ready, and message to notif player is error").Build())
			return
		}
	}

	stat, err := db.RiichiCommand.StartTournamentGame(match.TournamentId, playerReady, true)

	if err != nil && !stat {
		event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("some player not ready, and message to notif player is error").Build())
		return
	}

	event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("✅ Success Start the match").Build())
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
