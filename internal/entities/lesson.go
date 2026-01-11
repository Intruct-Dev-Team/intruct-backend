package entities

import "time"

type Lesson struct {
	ID           int       `db:"lesson_id"`
	CourseID     int       `db:"course_id"`
	ModuleID     int       `db:"module_id"`
	SerialNumber int       `db:"serial_number"`
	Title        string    `db:"title"`
	Description  string    `db:"description"`
	Content      string    `db:"content"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`

	Quizzes []*Quiz `db:"-"`
}
