package entities

import "time"

type Rating struct {
	ID        int       `db:"rating_id"`
	CourseID  int       `db:"course_id"`
	UserID    int       `db:"user_id"`
	Rating    int       `db:"rating"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
