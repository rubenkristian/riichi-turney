package discordbot

import (
	"context"
	"fmt"
	"log"
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
	Token    string
	AdminId  []string
	ServerId string
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

func (db *DiscordBot) StartBot(token string, serverId string) error {
	if db.IsRunning {
		return fmt.Errorf("discord bot already running")
	}

	db.Setting.Token = token
	db.Setting.ServerId = serverId

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

	defer client.Close(context.Background())

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

	if data.CommandName() == "register" {
		db.EventRegister(event)
	}

	if data.CommandName() == "start-table" {
		db.EventStartTable(event)
	}

	if data.CommandName() == "check-table" {
		db.EventCheckTable(event)
	}

	if data.CommandName() == "schedule-time" {
		db.EventCheckSchedule(event)
	}

	if data.CommandName() == "check-point" {
		db.EventCheckPoint(event)
	}
}

func (db *DiscordBot) onComponentInteract(event *events.ComponentInteractionCreate) {

}

func (db *DiscordBot) GetToken() string {
	return db.Setting.Token
}

func (db *DiscordBot) GetAdminId() []string {
	return db.Setting.AdminId
}
