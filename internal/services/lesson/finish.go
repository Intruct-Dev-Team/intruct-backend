package lesson

import (
	"context"
	"errors"
	"fmt"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/entities"
	serviceErrs "github.com/Intruct-Dev-Team/intruct-backend/internal/errors"
)

func (s *LessonService) FinishLesson(ctx context.Context, lessonID int, userID int) error {
	lesson, err := s.checkUserAccess(ctx, lessonID, userID)
	if err != nil {
		return err
	}

	// get all lesson IDs ordered by serial number
	lessonIDs, err := s.repo.GetAllLessonIDsByCourseIDOrderedBySerial(ctx, lesson.CourseID)
	if err != nil {
		return fmt.Errorf("failed to get lesson IDs: %w", err)
	}

	// find requested lesson index
	requestedLessonIndex := -1
	for i, id := range lessonIDs {
		if id == lessonID {
			requestedLessonIndex = i
			break
		}
	}

	if requestedLessonIndex == -1 {
		return serviceErrs.ErrorLessonNotFound
	}

	// get user's course progression
	progression, err := s.repo.GetCourseProgressionByUserAndCourse(ctx, userID, lesson.CourseID)
	if err != nil && !errors.Is(err, serviceErrs.ErrorSelectEmpty) {
		return fmt.Errorf("failed to get course progression: %w", err)
	}

	// if no progression exists
	if progression == nil {
		// user can only finish the first lesson
		if requestedLessonIndex != 0 {
			return serviceErrs.ErrorLessonNotReached
		}

		// create initial progression
		newProgression := &entities.CourseProgression{
			UserID:          userID,
			CourseID:        lesson.CourseID,
			CurrentLessonID: lessonIDs[0],
		}

		err = s.repo.CreateCourseProgression(ctx, newProgression)
		if err != nil {
			return fmt.Errorf("failed to create course progression: %w", err)
		}

		progression = newProgression
	}

	// find current lesson index
	currentLessonIndex := -1
	for i, id := range lessonIDs {
		if id == progression.CurrentLessonID {
			currentLessonIndex = i
			break
		}
	}

	if currentLessonIndex == -1 {
		return fmt.Errorf("current lesson not found in course")
	}

	// user can only finish the current lesson
	if requestedLessonIndex != currentLessonIndex {
		if requestedLessonIndex > currentLessonIndex {
			return serviceErrs.ErrorLessonNotReached
		}
		return serviceErrs.ErrorLessonFinished
	}

	// update progression
	newFinishedCount := progression.FinishedLessonsCount + 1
	var newCurrentLessonID int
	var isFinished bool

	// check if this is the last lesson
	if currentLessonIndex == len(lessonIDs)-1 {
		newCurrentLessonID = lessonID // keep current lesson as the last one
		isFinished = true
	} else {
		newCurrentLessonID = lessonIDs[currentLessonIndex+1]
		isFinished = false
	}

	err = s.repo.UpdateCourseProgression(ctx, userID, lesson.CourseID, newCurrentLessonID, newFinishedCount, isFinished)
	if err != nil {
		return fmt.Errorf("failed to update course progression: %w", err)
	}

	return nil
}
