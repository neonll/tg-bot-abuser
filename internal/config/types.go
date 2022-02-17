package config

type PenaltyType string

const (
	KICK PenaltyType = "kick"
	MUTE PenaltyType = "mute"
)

type Period string

const (
	MINUTE Period = "minute"
	HOUR   Period = "hour"
	DAY    Period = "day"
)

var AllPeriods = []Period{MINUTE, HOUR, DAY}
