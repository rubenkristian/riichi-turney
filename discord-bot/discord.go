package discordbot

import (
	"context"
	"fmt"
	"log"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/rubenkristian/riichi-turney/database"
	riichicommand "github.com/rubenkristian/riichi-turney/riichi-command"
)

type DiscordSetting struct {
	Token   string
	AdminId []string
}

type DiscordBot struct {
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

func (db *DiscordBot) StartBot(token string) error {
	if db.IsRunning {
		return fmt.Errorf("discord bot already running")
	}

	db.Setting.Token = token

	client, err := disgo.New(
		db.Setting.Token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuildMessages,
				gateway.IntentMessageContent,
			),
		),
		bot.WithEventListenerFunc(db.onMessageCreate),
	)

	if err != nil {
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

	return nil
}

func (db *DiscordBot) EndBot() error {
	if !db.IsRunning {
		return fmt.Errorf("discord bot not running")
	}

	db.Client.Close(context.Background())
	db.IsRunning = false

	return nil
}

func (db *DiscordBot) onMessageCreate(event *events.MessageCreate) {
	if event.Message.Author.Bot {
		return
	}

	fmt.Println(event.Message.Content)

	event.Client().Rest().AddReaction(event.ChannelID, event.MessageID, "✅")
}

func (db *DiscordBot) GetToken() string {
	return db.Setting.Token
}

func (db *DiscordBot) GetAdminId() []string {
	return db.Setting.AdminId
}
