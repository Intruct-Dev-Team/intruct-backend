package entities

import "time"

type Course struct {
	ID                 int        `db:"course_id"`
	OwnerID            int        `db:"owner_id"`
	Title              string     `db:"title"`
	Description        string     `db:"description"`
	LessonsCount       int        `db:"lessons_count"`
	LanguageID         int        `db:"language_id"`
	StateMachineItemID int        `db:"state_machine_item_id"`
	CreatedAt          time.Time  `db:"created_at"`
	UpdatedAt          time.Time  `db:"updated_at"`
	DeletedAt          *time.Time `db:"deleted_at"`
	IsPublic           bool       `db:"is_public"`

	Language          string             `db:"language"`
	State             string             `db:"state"`
	Modules           []*Module          `db:"-"`
	StateMachineName  StateMachineName   `db:"-"`
	StateMachineItem  *StateMachineItem  `db:"-"`
	CourseProgression *CourseProgression `db:"-"`
	Statistic         CourseStatistic    `db:"-"`
}

type CourseStatistic struct {
	CourseID      int     `db:"course_id"`
	StudentsCount int     `db:"students_count"`
	AverageRating float64 `db:"average_rating"`
	RatingsCount  int     `db:"ratings_count"`
}
