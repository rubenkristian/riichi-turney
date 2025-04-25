package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rubenkristian/riichi-turney/database"
	discordbot "github.com/rubenkristian/riichi-turney/discord-bot"
	riichicommand "github.com/rubenkristian/riichi-turney/riichi-command"
	"github.com/rubenkristian/riichi-turney/utils"
)

type ClientInterface struct {
	dbGame        *database.DatabaseGame
	riichiCommand *riichicommand.RiichiApi
	discordbot    *discordbot.DiscordBot
}

type DiscordSetting struct {
	Token string `json:"token"`
}

type CommandResponse struct {
}

func CreateClient(dbGame *database.DatabaseGame, riichiCommand *riichicommand.RiichiApi, discordbot *discordbot.DiscordBot) *ClientInterface {
	return &ClientInterface{
		dbGame:        dbGame,
		riichiCommand: riichiCommand,
		discordbot:    discordbot,
	}
}

func (ci *ClientInterface) SendCommand(c *fiber.Ctx) error {
	return utils.ResponseSuccess(200, "success", CommandResponse{})(c)
}

func (ci *ClientInterface) UpdateSetting(c *fiber.Ctx) error {
	return utils.ResponseSuccess(200, "success", CommandResponse{})(c)
}

func (ci *ClientInterface) StartRiichiBot(c *fiber.Ctx) error {
	return utils.ResponseSuccess(200, "success", CommandResponse{})(c)
}

func (ci *ClientInterface) EndRiichiBot(c *fiber.Ctx) error {
	return utils.ResponseSuccess(200, "success", CommandResponse{})(c)
}

func (ci *ClientInterface) StartDiscordBot(c *fiber.Ctx) error {
	discordSetting := new(DiscordSetting)

	if err := c.BodyParser(discordSetting); err != nil {
		return utils.ResponseError(fiber.StatusBadRequest, "Bad Request", err)(c)
	}

	if err := ci.discordbot.StartBot(discordSetting.Token); err != nil {
		return utils.ResponseError(fiber.StatusInternalServerError, "Some think went wrong", err)(c)
	}
	return utils.ResponseSuccess(200, "success", CommandResponse{})(c)
}

func (ci *ClientInterface) EndDiscordBot(c *fiber.Ctx) error {
	if err := ci.discordbot.EndBot(); err != nil {
		return utils.ResponseError(fiber.StatusInternalServerError, "Some think went wrong", err)(c)
	}

	return utils.ResponseSuccess(200, "success", CommandResponse{})(c)
}
