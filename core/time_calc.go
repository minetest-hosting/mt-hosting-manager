package core

import (
	"time"
)

// TODO: increment expiration time properly
func AddMonths(t time.Time, months int) time.Time {
	d := time.Hour * 24 * 31 * time.Duration(months)
	return t.Add(d)
}
