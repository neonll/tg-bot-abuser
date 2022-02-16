package main

import (
	"elisBot/internal/config"
	"elisBot/internal/stats"
	"elisBot/internal/tgbot"
)

func main() {
	config.ReadConfig()
	stats.Clear()
	tgbot.Run()
}
