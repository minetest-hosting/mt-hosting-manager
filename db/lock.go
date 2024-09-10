package db

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type DBLock struct {
	g *gorm.DB
}

func NewLock(g *gorm.DB) *DBLock {
	return &DBLock{g: g}
}

func (l *DBLock) TryLock(id int64) (bool, error) {
	result := false
	err := l.g.Raw("select * from pg_try_advisory_lock(?);", id).Scan(&result).Error
	return result, err
}

func (l *DBLock) UnLock(id int64) (bool, error) {
	result := false
	err := l.g.Raw("select * from pg_advisory_unlock(?);", id).Scan(&result).Error
	return result, err
}

type RunLockedFunc func() error

var ErrTryLockTimeout = fmt.Errorf("trylock timed out")

func (l *DBLock) RunLocked(id int64, timeout time.Duration, fn RunLockedFunc) error {
	start := time.Now()
	for {
		result, err := l.TryLock(id)
		if err != nil {
			return fmt.Errorf("trylock loop error: %v", err)
		}
		if result {
			break
		}
		time.Sleep(500 * time.Millisecond)
		if time.Since(start) > timeout {
			return ErrTryLockTimeout
		}
	}

	fnerr := fn()

	_, err := l.UnLock(id)
	if err != nil {
		return fmt.Errorf("unlock error: %v", err)
	}

	if fnerr != nil {
		return fmt.Errorf("user-function error: %v", err)
	}

	return nil
}
