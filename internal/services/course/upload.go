package course

import (
	"context"
	"errors"
	"fmt"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/entities"
	serviceErrs "github.com/Intruct-Dev-Team/intruct-backend/internal/errors"
)

// UploadCourseContent fills course with content from n8n processing
func (s *CourseService) UploadCourseContent(ctx context.Context, course *entities.Course) error {
	// Check course exists and get its current state
	existingCourse, err := s.repo.GetCourseByID(ctx, course.ID)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return serviceErrs.ErrorCourseNotFound
		}

		return fmt.Errorf("failed to get course: %w", err)
	}

	nextStateID, err := s.repo.GetStateIDByName(ctx, entities.Created)
	if err != nil {
		return fmt.Errorf("failed to get next stateID by name: %w", err)
	}

	// Get current state of course
	stateMachineItem, err := s.repo.GetStateMachineItemByID(ctx, existingCourse.StateMachineItemID)
	if err != nil {
		return fmt.Errorf("failed to get state machine item: %w", err)
	}

	available, err := s.repo.CheckIsTransitionAvailable(
		ctx,
		stateMachineItem.StateMachineID,
		stateMachineItem.StateID,
		nextStateID)
	if err != nil {
		return fmt.Errorf("failed to check if transition is available: %w", err)
	}

	if !available {
		return serviceErrs.ErrorUnavailableStateTransition
	}

	// Prepare course data for update
	course.Title = existingCourse.Title
	if course.Description == "" {
		course.Description = existingCourse.Description
	}
	course.StateMachineItem = stateMachineItem
	course.StateMachineItemID = existingCourse.StateMachineItemID
	course.OwnerID = existingCourse.OwnerID
	course.LanguageID = existingCourse.LanguageID

	if err := s.repo.ImplementCourse(ctx, course, nextStateID); err != nil {
		return fmt.Errorf("failed to implement course: %w", err)
	}

	return nil
}
