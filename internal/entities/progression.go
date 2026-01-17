package entities

type CourseProgression struct {
	UserID          int `db:"user_id"`
	CourseID        int `db:"course_id"`
	CurrentLessonID int `db:"current_lesson_id"`
}
