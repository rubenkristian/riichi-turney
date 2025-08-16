package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/rubenkristian/riichi-turney/database"
	discordbot "github.com/rubenkristian/riichi-turney/discord-bot"
	riichicommand "github.com/rubenkristian/riichi-turney/riichi-command"
)

type AppService struct {
	DbGame         *database.DatabaseGame
	RiichiCommand  *riichicommand.RiichiApi
	DiscordClient  bot.Client
	DiscordSetting *discordbot.DiscordSetting
}

func StartAppService(
	dbGame *database.DatabaseGame,
	riichiCommand *riichicommand.RiichiApi,
	discordClient *bot.Client,
	discordSetting *discordbot.DiscordSetting,
) *AppService {
	return &AppService{
		DbGame:         dbGame,
		RiichiCommand:  riichiCommand,
		DiscordClient:  *discordClient,
		DiscordSetting: discordSetting,
	}
}

func (as *AppService) SendInvite(tournamentId uint64) error {
	tournament, err := as.DbGame.GetTournament(tournamentId)

	if err != nil {
		return err
	}

	if !tournament.Active {
		return fmt.Errorf("Tournament is inactive")
	}

	players, err := as.DbGame.GetRegisterTournamentPlayers(tournamentId)
	fmt.Println(players)

	if err != nil {
		return err
	}

	if err := as.RiichiCommand.SendInvite(tournamentId, players); err != nil {
		return err
	}

	return nil
}

func (as *AppService) StartTable(tableId uint64) error {
	match, err := as.DbGame.GetMatchById(tableId)

	if err != nil {
		return fmt.Errorf("Match Not Found")
	}

	lobbyPlayers, err := as.RiichiCommand.FetchTournamentPlayers(match.TournamentId)

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
		_, err := as.DiscordClient.Rest().CreateMessage(
			snowflake.MustParse(as.DiscordSetting.ChannelNotify),
			discord.NewMessageCreateBuilder().SetContent(content).Build(),
		)

		if err != nil {
			return fmt.Errorf("some player not ready, and message to notif player is error")
		}
	}

	stat, err := as.RiichiCommand.StartTournamentGame(match.TournamentId, playerReady, true)

	if err != nil && !stat {
		return fmt.Errorf("some player not ready, and message to notif player is error")
	}

	return nil
}

func (as *AppService) FetchTournamentInfo(tournamentId uint64, registerEnd time.Time, roleID string) (*database.Tournament, error) {
	tournament, err := as.RiichiCommand.FetchTournamentInfo(tournamentId)

	if err != nil {
		return nil, err
	}

	newTournament, err := as.DbGame.CreateTournament(database.TournamentBody{
		Name:         tournament.MatchInfo.Name,
		Description:  tournament.MatchInfo.BriefIntroduction,
		StartAt:      time.Unix(tournament.MatchInfo.StartTime, 0),
		EndAt:        time.Unix(tournament.MatchInfo.EndTime, 0),
		RegisterEnd:  registerEnd,
		ClassifyID:   tournament.ClassifyID,
		RoleID:       roleID,
		TournamentId: tournamentId,
	})

	if err != nil {
		return nil, err
	}

	return newTournament, nil
}

func (as *AppService) FetchTournamentMatch(tournamentId uint64) error {
	tournament, err := as.DbGame.GetTournament(tournamentId)

	if err != nil {
		return err
	}

	lastId := 0

	for {
		tournamentMatches, err := as.RiichiCommand.FetchTournamentLogList(tournament.ClassifyID, lastId)

		if err != nil {
			return err
		}

		if len(tournamentMatches) == 0 {
			break
		}

		for _, tournamentMatch := range tournamentMatches {
			players := []database.InputTournamentMatchPlayer{}
			for _, player := range tournamentMatch.Players {
				players = append(players, database.InputTournamentMatchPlayer{
					PlayerId: player.UserId,
					Score:    player.Points,
				})
			}
			as.DbGame.CreateTournamentMatch(database.TournamentMatchBody{
				RoomID:       tournamentMatch.RoomID,
				PaiPuId:      tournamentMatch.PaiPuId,
				TournamentId: tournamentId,
				Players:      players,
			})
		}

		lastId += 20
	}

	return nil
}

func (as *AppService) FetchDetailTournamentMatch(tournamentMatchId string, paiPuId string) error {
	detailTournamentMatch, err := as.RiichiCommand.FetchLog(paiPuId)

	if err != nil {
		return err
	}

	lenHandRecord := len(detailTournamentMatch.HandRecords)
	lastHandRecord := detailTournamentMatch.HandRecords[lenHandRecord-1]
	lenHandEventRecord := len(lastHandRecord.HandEventRecords)
	lastHandEventRecord := lastHandRecord.HandEventRecords[lenHandEventRecord-1]

	var data DataHandEvent

	if err := json.Unmarshal([]byte(lastHandEventRecord.Data), &data); err != nil {
		return err
	}

	for _, userData := range data.UserData {
		if err := as.DbGame.UpdatePointTournamentMatchPlayer(tournamentMatchId, userData.UserID, userData.Score); err != nil {
			continue
		}
	}

	return nil
}

func (as *AppService) InputTournamentMatchPlayerPenalty(tournamentMatchPlayerId uint64, penalty int64) error {
	penalty = penalty * 10

	if err := as.DbGame.UpdateFinalPointTournamentMatchPlayer(tournamentMatchPlayerId, penalty); err != nil {
		return err
	}

	return nil
}

type DataHandEvent struct {
	UserData []struct {
		UserID   uint64 `json:"user_id"`
		PointNum int64  `json:"point_num"`
		Score    int64  `json:"score"`
	} `json:"user_data"`
}
