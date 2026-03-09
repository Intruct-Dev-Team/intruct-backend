package errors

import "errors"

var (
	ErrorUnavailableOperationState  = errors.New("unavailable operation for this state")
	ErrorUnavailableStateTransition = errors.New("unavailable such state transition")

	ErrorForbidden = errors.New("access is forbidden")

	ErrorUserExists   = errors.New("user already exists")
	ErrorUserNotFound = errors.New("user not found")

	ErrorLanguageNotFound           = errors.New("language not found")
	ErrorCourseExists               = errors.New("course already exists")
	ErrorCourseNotFound             = errors.New("course not found")
	ErrorCourseAlreadyPublished     = errors.New("course already published")
	ErrorNotCourseOwner             = errors.New("user is not the owner of the course")
	ErrorUserHasNoCourseProgression = errors.New("user have not learn this course")

	ErrorLessonNotFound   = errors.New("lesson not found")
	ErrorLessonNotReached = errors.New("previous lessons were not finished")
	ErrorLessonFinished   = errors.New("lesson is already finished")

	ErrorRatingExists = errors.New("user already left his grade")
)
