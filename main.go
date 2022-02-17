package main

import (
	"elisBot/internal/config"
	"elisBot/internal/stats"
	"elisBot/internal/tgbot"
)

func main() {
	config.ReadConfig()
	go stats.ScheduleClear()
	tgbot.Run()
}
