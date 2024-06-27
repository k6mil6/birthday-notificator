package offset

import "time"

type Offset struct {
	Unit  string `json:"unit"`
	Value int    `json:"value"`
}

func ConvertToTimeDuration(offset Offset) time.Duration {
	switch offset.Unit {
	case "day":
		return time.Duration(offset.Value) * time.Hour * 24
	case "hour":
		return time.Duration(offset.Value) * time.Hour
	case "minute":
		return time.Duration(offset.Value) * time.Minute
	case "second":
		return time.Duration(offset.Value) * time.Second
	default:
		return 0
	}
}
