package main

import (
	"embed"
	_ "embed"
	"io/fs"
	"log/slog"
	"net/http"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/rubenkristian/riichi-turney/database"
	discordbot "github.com/rubenkristian/riichi-turney/discord-bot"
	riichicommand "github.com/rubenkristian/riichi-turney/riichi-command"
	"github.com/rubenkristian/riichi-turney/server"
)

//go:embed client/dist/*
var webApp embed.FS

func serveStaticFile(c *fiber.Ctx) error {
	path := "client/dist" + c.Path()

	if path == "client/dist/" {
		path = "client/dist/index.html"
	}

	slog.Info(path)

	file, err := webApp.Open(path)
	if err != nil {
		return c.Status(http.StatusNotFound).SendString("File not found")
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	data, _ := fs.ReadFile(webApp, path)

	ext := filepath.Ext(path)
	contentType := getContentType(ext)

	// Set the Last-Modified header
	c.Set("Last-Modified", fileInfo.ModTime().UTC().Format(http.TimeFormat))

	// Set the Content-Type header
	c.Set("Content-Type", contentType)
	// Return the file content
	return c.Send(data)
}

func getContentType(ext string) string {
	switch ext {
	case ".html":
		return "text/html; charset=utf-8"
	case ".css":
		return "text/css; charset=utf-8"
	case ".js":
		return "application/javascript; charset=utf-8"
	case ".json":
		return "application/json"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	default:
		return "application/octet-stream"
	}
}

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
	app.Get("*", serveStaticFile)

	app.Post("/discord/power/on", client.StartDiscordBot)
	app.Post("/discord/power/off", client.EndDiscordBot)
	app.Listen(":8080")
}
