package main

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/rubenkristian/riichi-turney/database"
	discordbot "github.com/rubenkristian/riichi-turney/discord-bot"
	riichicommand "github.com/rubenkristian/riichi-turney/riichi-command"
	"github.com/rubenkristian/riichi-turney/server"
)

func main() {
	slog.Info("Riichi turney")
	dbGame := database.CreateDatabaseGame()
	riichiComand := riichicommand.CreateRiichiApi(dbGame)
	discordBot := discordbot.CreateDiscordBot(
		dbGame,
		riichiComand,
	)

	client := server.CreateClient(dbGame, riichiComand, discordBot)

	app := fiber.New()

	app.Post("/discord/power/on", client.StartDiscordBot)
	app.Post("/discord/power/off", client.EndDiscordBot)
	app.Listen(":8080")
}
