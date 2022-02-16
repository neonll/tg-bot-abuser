package penalties

import (
	"elisBot/internal/config"
	"elisBot/internal/stats"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"time"
)

func Fire(action string, update tgbotapi.Update) []tgbotapi.Chattable {
	switch action {
	case "kick":
		return check(config.KICK, update)
	case "mute":
		return check(config.MUTE, update)
	default:
		return nil
	}
}

func check(action config.PenaltyType, update tgbotapi.Update) []tgbotapi.Chattable {
	result := make([]tgbotapi.Chattable, 0)
	hour := update.Message.Time().Hour()
	minute := update.Message.Time().Minute()

	for _, period := range config.All {
		count := stats.GetCount(period, hour, minute)
		limit := config.Limits[action][period]
		if count >= limit {
			return fire(action, period, update)
		} else {
			left := limit - count
			threshold := config.Thresholds[action][period]
			if left <= threshold {
				return warn(action, period, left, update)
			}
		}
	}
	return result
}

func fire(action config.PenaltyType, period config.Period, update tgbotapi.Update) []tgbotapi.Chattable {
	result := make([]tgbotapi.Chattable, 0)

	log.Printf("%v for %d minutes!", action, config.Penalties[action][period])
	result = append(result,
		replyMessageFire[action](update, config.Penalties[action][period]),
		actionFire[action](update, config.Penalties[action][period]))

	return result
}

func warn(action config.PenaltyType, period config.Period, left int, update tgbotapi.Update) []tgbotapi.Chattable {
	result := make([]tgbotapi.Chattable, 0)
	log.Printf("Warning! Left %d messages before %v for %d minutes!", left, action, config.Penalties[action][period])
	msg := replyMessageWarn[action](update, left, config.Penalties[action][period])
	return append(result, msg)
}

var actionFire = map[config.PenaltyType]func(tgbotapi.Update, int) tgbotapi.Chattable{
	config.MUTE: func(update tgbotapi.Update, minutes int) tgbotapi.Chattable {
		return tgbotapi.RestrictChatMemberConfig{
			ChatMemberConfig: getChatMemberConfig(update),
			UntilDate:        getUntilDate(update, minutes),
			Permissions:      &config.RestrictAll,
		}
	},
	config.KICK: func(update tgbotapi.Update, minutes int) tgbotapi.Chattable {
		return tgbotapi.BanChatMemberConfig{
			ChatMemberConfig: getChatMemberConfig(update),
			UntilDate:        getUntilDate(update, minutes),
			RevokeMessages:   false,
		}
	},
}

func getChatMemberConfig(update tgbotapi.Update) tgbotapi.ChatMemberConfig {
	return tgbotapi.ChatMemberConfig{
		ChatID: update.FromChat().ID,
		UserID: update.SentFrom().ID,
	}
}

func getUntilDate(update tgbotapi.Update, minutes int) int64 {
	return update.Message.Time().Add(time.Minute * time.Duration(minutes)).Unix()
}

var replyMessageFire = map[config.PenaltyType]func(tgbotapi.Update, int) tgbotapi.Chattable{
	config.MUTE: func(update tgbotapi.Update, minutes int) tgbotapi.Chattable {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			fmt.Sprintf("Договорился! Молчанка на %d минут!", minutes))
		msg.ReplyToMessageID = update.Message.MessageID
		return msg
	},
	config.KICK: func(update tgbotapi.Update, minutes int) tgbotapi.Chattable {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			fmt.Sprintf("Договорился! Бан на %d минут!", minutes))
		msg.ReplyToMessageID = update.Message.MessageID
		return msg
	},
}

var replyMessageWarn = map[config.PenaltyType]func(tgbotapi.Update, int, int) tgbotapi.Chattable{
	config.MUTE: func(update tgbotapi.Update, left, minutes int) tgbotapi.Chattable {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			fmt.Sprintf("Притормози! Через %d сообщений будет молчанка на %d минут!", left, minutes))
		msg.ReplyToMessageID = update.Message.MessageID
		return msg
	},
	config.KICK: func(update tgbotapi.Update, left, minutes int) tgbotapi.Chattable {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			fmt.Sprintf("Притормози! Через %d сообщений будет бан на %d минут!", left, minutes))
		msg.ReplyToMessageID = update.Message.MessageID
		return msg
	},
}
