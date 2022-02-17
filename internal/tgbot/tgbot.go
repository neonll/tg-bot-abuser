package tgbot

import (
	"elisBot/internal/config"
	"elisBot/internal/penalties"
	"elisBot/internal/stats"
	. "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Run() {
	bot, err := NewBotAPI(config.BotToken)
	if err != nil {
		panic(err)
	}

	bot.Debug = config.BotDebug

	updateConfig := NewUpdate(0)

	updateConfig.Timeout = config.BotTimeout

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		// Telegram can send many types of updates depending on what your Bot
		// is up to. We only want to look at messages for now, so we can
		// discard any other updates.
		if update.Message == nil {
			continue
		}

		doSomething(update, *bot)
	}
}

func doSomething(update Update, bot BotAPI) {

	if !int64SliceContains(config.ManagedChats, update.FromChat().ID) {
		return
	}

	if config.ManagedUser != update.SentFrom().ID {
		return
	}

	hour := update.Message.Time().Hour()
	minute := update.Message.Time().Minute()
	stats.Inc(hour, minute)

	for _, action := range config.ActivePenalties {
		results := penalties.Fire(action, update)
		for _, result := range results {
			switch result.(type) {
			case MessageConfig:
				if _, err := bot.Send(result); err != nil {
					panic(err)
				}
			case RestrictChatMemberConfig:
				if _, err := bot.Request(result); err != nil {
					panic(err)
				}
			case BanChatMemberConfig:
				if _, err := bot.Request(result); err != nil {
					panic(err)
				}
			}
		}
	}
}

func int64SliceContains(s []int64, e int64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
