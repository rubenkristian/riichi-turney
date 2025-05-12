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
	Token    string `json:"token"`
	ServerId string `json:"server"`
}

type RiichiCitySetting struct {
	Domain   string `json:"domain"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type CommandResponse struct {
}

type DiscordStatusResponse struct {
	Status    bool  `json:"status"`
	StartTime int64 `json:"start_time"`
}

type CreateResponse struct {
	Status bool `json:"status"`
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
	riichiCitySetting := new(RiichiCitySetting)

	if err := c.BodyParser(riichiCitySetting); err != nil {
		return utils.ResponseError(fiber.StatusBadRequest, "Bad Request", err)(c)
	}

	if err := ci.riichiCommand.SetupRiichi(riichiCitySetting.Domain, riichiCitySetting.Username, riichiCitySetting.Password); err != nil {
		return utils.ResponseError(fiber.StatusInternalServerError, "Some think went wrong", err)(c)
	}

	return utils.ResponseSuccess(200, "success", CommandResponse{})(c)
}

func (ci *ClientInterface) EndRiichiBot(c *fiber.Ctx) error {
	if !ci.riichiCommand.IsLoggedIn {
		return utils.ResponseError(fiber.StatusInternalServerError, "Riichi Bot already end", nil)(c)
	}

	return utils.ResponseSuccess(200, "success", CommandResponse{})(c)
}

func (ci *ClientInterface) StartDiscordBot(c *fiber.Ctx) error {
	discordSetting := new(DiscordSetting)

	if err := c.BodyParser(discordSetting); err != nil {
		return utils.ResponseError(fiber.StatusBadRequest, "Bad Request", err)(c)
	}

	if err := ci.discordbot.StartBot(discordSetting.Token, discordSetting.ServerId); err != nil {
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

func (ci *ClientInterface) CheckStatusDiscordBot(c *fiber.Ctx) error {
	return utils.ResponseSuccess(fiber.StatusOK, "success", DiscordStatusResponse{
		Status:    ci.discordbot.IsRunning,
		StartTime: ci.discordbot.StartTime,
	})(c)
}

func (ci *ClientInterface) CreatePlayer(c *fiber.Ctx) error {
	return utils.ResponseSuccess(fiber.StatusOK, "success", CreateResponse{})(c)
}

func (ci *ClientInterface) UpdatePlayer(c *fiber.Ctx) error {
	return utils.ResponseSuccess(fiber.StatusOK, "success", CreateResponse{})(c)
}

func (ci *ClientInterface) DeletePlayer(c *fiber.Ctx) error {
	return utils.ResponseSuccess(fiber.StatusOK, "success", CreateResponse{})(c)
}

func (ci *ClientInterface) ViewPlayer(c *fiber.Ctx) error {
	return utils.ResponseSuccess(fiber.StatusOK, "success", CreateResponse{})(c)
}

func (ci *ClientInterface) CreateHanchan(c *fiber.Ctx) error {
	return utils.ResponseSuccess(fiber.StatusOK, "success", CreateResponse{})(c)
}

func (ci *ClientInterface) CreateTable(c *fiber.Ctx) error {
	return utils.ResponseSuccess(fiber.StatusOK, "success", CreateResponse{})(c)
}

func (ci *ClientInterface) UpdateTable(c *fiber.Ctx) error {
	return utils.ResponseSuccess(fiber.StatusOK, "success", CreateResponse{})(c)
}

func (ci *ClientInterface) StartTable(c *fiber.Ctx) error {
	return utils.ResponseSuccess(fiber.StatusOK, "success", CreateResponse{})(c)
}

func (ci *ClientInterface) CheckTable(c *fiber.Ctx) error {
	return utils.ResponseSuccess(fiber.StatusOK, "success", CreateResponse{})(c)
}
