package entities

import "time"

type Quiz struct {
	ID           int       `db:"quiz_id"`
	LessonID     int       `db:"lesson_id"`
	SerialNumber int       `db:"serial_number"`
	Question     string    `db:"question"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`

	Options []*QuizOption `db:"-"`
}

type QuizOption struct {
	ID        int       `db:"option_id"`
	QuizID    int       `db:"quiz_id"`
	Content   string    `db:"content"`
	IsAnswer  bool      `db:"is_answer"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
