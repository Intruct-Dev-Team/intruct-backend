package entities

import "time"

type Streak struct {
	UserID              int       `db:"user_id"`
	DaysStreakCount     int       `db:"days_streak_count"`
	CreatedAt           time.Time `db:"created_at"`
	UpdatedAt           time.Time `db:"updated_at"`
	IsStreakActiveToday bool      `db:"-"`
}
