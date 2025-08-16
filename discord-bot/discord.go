package discordbot

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
	"github.com/rubenkristian/riichi-turney/database"
	riichicommand "github.com/rubenkristian/riichi-turney/riichi-command"
)

type DiscordSetting struct {
	Token         string
	AdminId       []string
	ServerId      string
	ChannelAdmin  string
	ChannelNotify string
}

type DiscordBot struct {
	StartTime     int64
	Setting       DiscordSetting
	Client        bot.Client
	DbGame        *database.DatabaseGame
	RiichiCommand *riichicommand.RiichiApi
	IsRunning     bool
}

func CreateDiscordBot(dbGame *database.DatabaseGame, riichiCommand *riichicommand.RiichiApi) *DiscordBot {
	return &DiscordBot{
		DbGame:        dbGame,
		RiichiCommand: riichiCommand,
		IsRunning:     false,
	}
}

func (db *DiscordBot) StartBot(token string, serverId string, channelAdmin string, channelNotify string) error {
	if db.IsRunning {
		return fmt.Errorf("discord bot already running")
	}

	db.Setting.Token = token
	db.Setting.ServerId = serverId
	db.Setting.ChannelAdmin = channelAdmin
	db.Setting.ChannelNotify = channelNotify

	client, err := disgo.New(
		db.Setting.Token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuildMessages,
				gateway.IntentMessageContent,
			),
		),
	)

	client.AddEventListeners(bot.NewListenerFunc(db.onMessageInteract))
	client.AddEventListeners(bot.NewListenerFunc(db.onEventInteract))
	client.AddEventListeners(bot.NewListenerFunc(db.onComponentInteract))

	if err != nil {
		return err
	}

	commands := []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:        "register",
			Description: "register turney",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "riichi_city_id",
					Description: "Id user riichi city",
					Required:    true,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "start-table",
			Description: "start table with number of table",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "table_id",
					Description: "id of table active",
					Required:    true,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "check-table",
			Description: "check table (all or one)",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "table_id",
					Description: "id of table active",
					Required:    false,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "schedule-time",
			Description: "get detail schedule time",
		},
		discord.SlashCommandCreate{
			Name:        "check-point",
			Description: "check point of player for current turney",
		},
	}

	if _, err := client.Rest().SetGuildCommands(client.ApplicationID(), snowflake.MustParse(db.Setting.ServerId), commands); err != nil {
		return err
	}

	go func() {
		if err = client.OpenGateway(context.Background()); err != nil {
			log.Printf("error opening gateway: %v", err)
		}
	}()

	db.Client = client
	db.IsRunning = true
	db.StartTime = time.Now().UnixMilli()

	return nil
}

func (db *DiscordBot) EndBot() error {
	if !db.IsRunning {
		return fmt.Errorf("discord bot not running")
	}

	db.Client.Close(context.Background())
	db.IsRunning = false
	db.StartTime = 0

	return nil
}

func (db *DiscordBot) onMessageInteract(event *events.MessageCreate) {
	if event.Message.Author.Bot {
		return
	}

	event.Client().Rest().AddReaction(event.ChannelID, event.MessageID, "✅")
}

func (db *DiscordBot) onEventInteract(event *events.ApplicationCommandInteractionCreate) {
	data := event.SlashCommandInteractionData()

	switch data.CommandName() {
	case "start-table":
		db.EventStartTable(event)

	case "check-table":
		db.EventCheckTable(event)

	case "register":
		db.EventRegister(event)

	case "schedule-time":
		db.EventCheckSchedule(event)

	case "check-point":
		db.EventCheckPoint(event)
	}
}

func (db *DiscordBot) onComponentInteract(event *events.ComponentInteractionCreate) {
	customID := event.Data.CustomID()

	parts := strings.SplitN(customID, ":", 2)
	action := parts[0]

	fmt.Println(customID)

	switch action {
	case "confirm":
		playerId, err := strconv.ParseUint(parts[1], 10, 64)
		if err != nil {
			event.UpdateMessage(
				discord.NewMessageUpdateBuilder().SetContent(fmt.Sprintf("❌ Registration cancelled, %s", err.Error())).
					SetContainerComponents().
					Build(),
			)
			return
		}

		newRegister, err := db.DbGame.CreateRegisterTournament(playerId)

		if err != nil {
			event.UpdateMessage(
				discord.NewMessageUpdateBuilder().SetContent(fmt.Sprintf("❌ Registration cancelled, %s", err.Error())).
					SetContainerComponents().
					Build(),
			)
			return
		}

		db.Client.Rest().AddMemberRole(snowflake.MustParse(db.Setting.ServerId), event.User().ID, snowflake.MustParse(newRegister.Tournament.RoleID))
		event.UpdateMessage(
			discord.NewMessageUpdateBuilder().
				SetContent(fmt.Sprintf("✅ Registration success for %s.", newRegister.Player.RiichiCityName)).
				SetContainerComponents().
				Build(),
		)
	case "cancel":
		event.UpdateMessage(
			discord.NewMessageUpdateBuilder().SetContent("❌ Registration cancelled").
				SetContainerComponents().
				Build(),
		)

	default:
		event.UpdateMessage(
			discord.NewMessageUpdateBuilder().SetContent("❌ no action for " + action).
				SetContainerComponents().
				Build(),
		)
	}
}

func (db *DiscordBot) GetToken() string {
	return db.Setting.Token
}

func (db *DiscordBot) GetAdminId() []string {
	return db.Setting.AdminId
}
