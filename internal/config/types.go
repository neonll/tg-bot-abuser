package config

type PenaltyType string

const (
	KICK PenaltyType = "kick"
	MUTE PenaltyType = "mute"
)

var AllPenaltyTypes = []PenaltyType{KICK, MUTE}
