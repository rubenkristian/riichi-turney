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
		event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Please to insert riichi id as number").SetEphemeral(true).Build())
		return
	}

	playerGame, err := db.DbGame.GetPlayerByRiichiId(riichiId)

	if err != nil {
		fmt.Println(err.Error())
	}

	player, err := db.RiichiCommand.FindPlayer(riichi_city_id)

	if err != nil {
		fmt.Println(err.Error())
		event.CreateMessage(discord.NewMessageCreateBuilder().SetEphemeral(true).SetContent(fmt.Sprintf("This id %s not found", riichi_city_id)).Build())
		return
	}

	if len(player.FriendList) <= 0 {
		event.CreateMessage(discord.NewMessageCreateBuilder().SetEphemeral(true).SetContent(fmt.Sprintf("This id %s not found", riichi_city_id)).Build())
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
			SetEphemeral(true).
			AddActionRow(confirmBtn, cancelBtn).Build(),
	)
}

// handle when command /start-table selected
func (db *DiscordBot) EventStartTable(event *events.ApplicationCommandInteractionCreate) {
	if uint64(event.Channel().ID()) != uint64(snowflake.MustParse(db.Setting.ChannelAdmin)) {
		event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Not Allowed").SetEphemeral(true).Build())
		return
	}

	data := event.SlashCommandInteractionData()

	match_id := data.String("table_id")

	matchId, err := strconv.ParseUint(match_id, 10, 64)

	if err != nil {
		event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Please to insert match id as number").SetEphemeral(true).Build())
		return
	}

	match, err := db.DbGame.GetMatchById(matchId)

	if err != nil {
		event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Match Not Found").SetEphemeral(true).Build())
		return
	}

	if match.Status == 1 {
		event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("This match has started").SetEphemeral(true).Build())
		return
	}

	lobbyPlayers, err := db.RiichiCommand.FetchTournamentPlayers(match.TournamentId)

	var playersMatch map[uint64]bool = map[uint64]bool{}
	var playerReady []uint64 = []uint64{}
	var playerNotReady []database.PlayerCheck = []database.PlayerCheck{}

	for _, player := range match.PlayerMatches {
		playersMatch[player.PlayerId] = false
	}

	for _, lobbyPlayer := range lobbyPlayers {
		if _, ok := playersMatch[lobbyPlayer.UserID]; ok && lobbyPlayer.Status == 2 {
			playersMatch[lobbyPlayer.UserID] = true
		}
	}

	for _, playerMatch := range match.PlayerMatches {
		if playersMatch[playerMatch.PlayerId] {
			playerReady = append(playerReady, playerMatch.PlayerId)
		} else {
			playerNotReady = append(playerNotReady, database.PlayerCheck{
				RiichiId:  playerMatch.PlayerId,
				DiscordId: playerMatch.Player.DiscordId,
				Status:    -1,
			})
		}
	}

	if len(playerNotReady) > 0 {
		mentions := ""

		for _, player := range playerNotReady {
			mentions += fmt.Sprintf("<@%d>\n", player.DiscordId)
		}
		content := "The following players, please get ready in the tournament lobby — the match is about to start!\n\n" + mentions
		_, err := db.Client.Rest().CreateMessage(
			snowflake.MustParse(db.Setting.ChannelNotify),
			discord.NewMessageCreateBuilder().SetContent(content).SetEphemeral(true).Build(),
		)

		if err != nil {
			event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("some player not ready, and message to notif player is error").SetEphemeral(true).Build())
			db.DbGame.UpdateStatusMatch(matchId, -1)
			return
		}
		event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("❌ Failed Start the match").SetEphemeral(true).Build())

		return
	}

	stat, err := db.RiichiCommand.StartTournamentGame(match.TournamentId, playerReady, true)

	fmt.Println(stat)

	if err != nil {
		event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("some player not ready, and message to notif player is error").SetEphemeral(true).Build())
		return
	}

	event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("✅ Success Start the match").SetEphemeral(true).Build())

	db.DbGame.UpdateStatusMatch(matchId, 1)
}

// handle when command /check-table selected
func (db *DiscordBot) EventCheckTable(event *events.ApplicationCommandInteractionCreate) {
	if uint64(event.Channel().ID()) != uint64(snowflake.MustParse(db.Setting.ChannelAdmin)) {
		event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Not Allowed").SetEphemeral(true).Build())
		return
	}
	matches, err := db.DbGame.ListNotStartedMatch()

	if err != nil {
		event.CreateMessage(discord.NewMessageCreateBuilder().SetContent(err.Error()).Build())
		return
	}

	list := "\nPending Matches:\n"

	for _, match := range matches {
		players := ""

		for id, player := range match.PlayerMatches {
			players += fmt.Sprintf("- Player %d : <@%d>\n", id+1, player.Player.DiscordId)
		}
		fmt.Println(players)
		list += fmt.Sprintf("**Match (id: %d ):**\n %s", match.Id, players)
	}

	event.CreateMessage(discord.NewMessageCreateBuilder().SetContent(list).Build())
}

// handle when command /check-schedule selected
func (db *DiscordBot) EventCheckSchedule(event *events.ApplicationCommandInteractionCreate) {

}

// handle when command /check-point selected
func (db *DiscordBot) EventCheckPoint(event *events.ApplicationCommandInteractionCreate) {
	user := event.User()
	discordID := uint64(user.ID)

	player, err := db.DbGame.GetPlayerByDiscordId(discordID)

	if err != nil {
		fmt.Println(err)
		event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("You riichi id not found").Build())
		return
	}

	points, err := db.DbGame.PointsByPlayer(player.Id)

	if err != nil {
		fmt.Println(err)
		event.CreateMessage(discord.NewMessageCreateBuilder().SetContent(err.Error()).Build())
		return
	}

	list := fmt.Sprintf("\n**Points** for <@%d>\n", discordID)

	for _, point := range points {
		list += fmt.Sprintf("point: %d | penalty: %d | final point: %d\n", point.Point, point.Penalty, point.FinalPoint)
	}

	fmt.Println(points)

	event.CreateMessage(discord.NewMessageCreateBuilder().SetContent(list).Build())
}
