package entities

type CourseProgression struct {
	UserID               int  `db:"user_id"`
	CourseID             int  `db:"course_id"`
	CurrentLessonID      int  `db:"current_lesson_id"`
	FinishedLessonsCount int  `db:"finished_lessons_count"`
	IsFinished           bool `db:"is_finished"`
}
