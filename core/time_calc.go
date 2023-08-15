package core

import (
	"time"
)

// TODO: increment expiration time properly
func AddDays(t time.Time, days int) time.Time {
	d := time.Hour * 24 * time.Duration(days)
	return t.Add(d)
}
