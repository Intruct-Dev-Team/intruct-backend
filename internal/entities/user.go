package entities

import "time"

type User struct {
	ID               int       `db:"user_id"`
	ExternalUUID     string    `db:"external_uuid"`
	Email            string    `db:"email"`
	Name             string    `db:"name"`
	Surname          string    `db:"surname"`
	RegistrationDate time.Time `db:"registration_date"`
	Birthdate        time.Time `db:"birthdate"`
	Avatar           string    `db:"avatar"`

	Statistic UserStatistic `db:"-"`
	Streak    Streak        `db:"-"`
}

type UserStatistic struct {
	CountOfCompletedCourses  int `db:"count_of_complited_courses"`
	CountOfInProgressCourses int `db:"count_of_in_progress_courses"`
}
