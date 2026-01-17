package course

import (
	"context"
	"errors"
	"fmt"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/entities"
	serviceErrs "github.com/Intruct-Dev-Team/intruct-backend/internal/errors"
)

// PublishCourse set course in published state
func (s *CourseService) PublishCourse(ctx context.Context, courseID int, userID int) error {
	exists, err := s.repo.IsUserExistsByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check user existence by id: %w", err)
	}
	if !exists {
		return serviceErrs.ErrorUserNotFound
	}

	// Check course exists and get its current state
	course, err := s.repo.GetCourseByID(ctx, courseID)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return serviceErrs.ErrorCourseNotFound
		}

		return fmt.Errorf("failed to get course: %w", err)
	}

	if course.OwnerID != userID {
		return serviceErrs.ErrorNotCourseOwner
	}

	nextStateID, err := s.repo.GetStateIDByName(ctx, entities.Published)
	if err != nil {
		return fmt.Errorf("failed to get next stateID by name: %w", err)
	}

	// Get current state of course
	stateMachineItem, err := s.repo.GetStateMachineItemByID(ctx, course.StateMachineItemID)
	if err != nil {
		return fmt.Errorf("failed to get state machine item: %w", err)
	}

	if stateMachineItem.StateID == nextStateID {
		return serviceErrs.ErrorCourseAlreadyPublished
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

	if err := s.repo.UpdateStateMachineItemState(ctx, course.StateMachineItemID, nextStateID); err != nil {
		return fmt.Errorf("failed to publish course: %w", err)
	}

	if err := s.repo.SetIsPublicField(ctx, courseID, true); err != nil {
		return fmt.Errorf("failed to set is_public field for course: %w", err)
	}

	return nil
}
