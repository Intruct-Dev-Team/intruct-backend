package entities

import "time"

type Course struct {
	ID                 int       `db:"course_id"`
	OwnerID            int       `db:"owner_id"`
	Title              string    `db:"title"`
	Description        string    `db:"description"`
	LessonsCount       int       `db:"lessons_count"`
	LanguageID         int       `db:"language_id"`
	StateMachineItemID int       `db:"state_machine_item_id"`
	CreatedAt          time.Time `db:"created_at"`
	UpdatedAt          time.Time `db:"updated_at"`
	IsPublic           bool      `db:"is_public"`

	Language          string             `db:"language"`
	Modules           []*Module          `db:"-"`
	StateMachineName  StateMachineName   `db:"-"`
	StateMachineItem  *StateMachineItem  `db:"-"`
	CourseProgression *CourseProgression `db:"-"`
}
