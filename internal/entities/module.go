package entities

import "time"

type Module struct {
	ID           int       `db:"module_id"`
	CourseID     int       `db:"course_id"`
	SerialNumber int       `db:"serial_number"`
	Title        string    `db:"title"`
	Description  string    `db:"description"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`

	Lessons []*Lesson `db:"-"`
}
