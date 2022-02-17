package config

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

var (
	configPath      = "config.yaml"
	BotToken        string
	BotDebug        bool
	BotTimeout      int
	ActivePenalties []string
	ManagedChats    []int64
	ManagedUser     int64
	Penalties       map[PenaltyType]map[Period]int
	Thresholds      map[PenaltyType]map[Period]int
	Limits          map[PenaltyType]map[Period]int
	RestrictAll     = tgbotapi.ChatPermissions{
		CanSendMessages:       false,
		CanSendMediaMessages:  false,
		CanSendPolls:          false,
		CanSendOtherMessages:  false,
		CanAddWebPagePreviews: false,
		CanChangeInfo:         false,
		CanInviteUsers:        false,
		CanPinMessages:        false,
	}
)

func ReadConfig() {
	type schema struct {
		Bot struct {
			Token   string `yaml:"token"`
			Debug   bool   `yaml:"debug"`
			Timeout int    `yaml:"timeout"`
		} `yaml:"bot"`
		ActivePenalties []string `yaml:"activePenalties,flow"`
		Managed         struct {
			Chats []int64 `yaml:"chats,flow"`
			User  int64   `yaml:"user"`
		} `yaml:"managed"`
		Limits     map[PenaltyType]map[Period]int `yaml:"limits"`
		Thresholds map[PenaltyType]map[Period]int `yaml:"thresholds"`
		Penalties  map[PenaltyType]map[Period]int `yaml:"penalties"`
	}

	t := schema{}

	configBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = yaml.Unmarshal(configBytes, &t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	BotToken = t.Bot.Token
	BotDebug = t.Bot.Debug
	BotTimeout = t.Bot.Timeout

	ActivePenalties = t.ActivePenalties

	ManagedChats = t.Managed.Chats
	ManagedUser = t.Managed.User

	Limits = t.Limits
	Thresholds = t.Thresholds
	Penalties = t.Penalties

}
