package config

type Period string

const (
	MINUTE Period = "minute"
	HOUR   Period = "hour"
	DAY    Period = "day"
)

var All = []Period{MINUTE, HOUR, DAY}
