package lesson

import (
	"context"
	"errors"
	"fmt"
	"time"

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
		err = s.repo.CreateCourseProgression(ctx, userID, lesson.CourseID, lessonIDs[0])
		if err != nil {
			return fmt.Errorf("failed to create course progression: %w", err)
		}

		progression = &entities.CourseProgression{
			UserID:          userID,
			CourseID:        lesson.CourseID,
			CurrentLessonID: &lessonIDs[0],
		}
	}

	// find current lesson index
	if progression.CurrentLessonID == nil {
		return fmt.Errorf("current lesson not set")
	}

	currentLessonIndex := -1
	for i, id := range lessonIDs {
		if id == *progression.CurrentLessonID {
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

	// calculate new progression values
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

	// get or create user streak
	streak, err := s.repo.GetUserStreakByUserID(ctx, userID)
	if err != nil && !errors.Is(err, serviceErrs.ErrorSelectEmpty) {
		return fmt.Errorf("failed to get user streak: %w", err)
	}

	// calculate streak values
	var newDaysStreak int

	if streak == nil {
		// no streak exists - will be created with initial values
		newDaysStreak = 1
	} else {
		now := time.Now().UTC()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

		updatedDay := time.Date(
			streak.UpdatedAt.UTC().Year(),
			streak.UpdatedAt.UTC().Month(),
			streak.UpdatedAt.UTC().Day(),
			0, 0, 0, 0,
			time.UTC,
		)

		if updatedDay.Equal(today) {
			// already completed today - no change
			newDaysStreak = streak.DaysStreakCount
		} else if updatedDay.Equal(today.AddDate(0, 0, -1)) {
			// completed yesterday - increment streak
			newDaysStreak = streak.DaysStreakCount + 1
		} else {
			// streak broken - reset
			newDaysStreak = 1
		}
	}

	// transactionally update progression and streak (or create streak if needed)
	err = s.repo.UpdateCourseProgressionAndStreak(
		ctx,
		userID,
		lesson.CourseID,
		newCurrentLessonID,
		newFinishedCount,
		isFinished,
		newDaysStreak,
	)
	if err != nil {
		return fmt.Errorf("failed to update course progression and streak: %w", err)
	}

	return nil
}
