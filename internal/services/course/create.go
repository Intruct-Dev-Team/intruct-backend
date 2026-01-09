package course

import (
	"context"
	"errors"
	"fmt"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/entities"
	serviceErrs "github.com/Intruct-Dev-Team/intruct-backend/internal/errors"
)

// CreateCourse creates new course, returns course id from db.
func (s *CourseService) CreateCourse(ctx context.Context, course *entities.Course) (int, error) {
	exists, err := s.repo.IsUserExistsByID(ctx, course.OwnerID)
	if err != nil {
		return 0, fmt.Errorf("failed to check user existence by id: %w", err)
	}
	if !exists {
		return 0, serviceErrs.ErrorUserNotFound
	}

	// get language ID by name
	languageID, err := s.repo.GetLanguageIDByName(ctx, course.Language)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return 0, fmt.Errorf("language '%s' not found: %w", course.Language, serviceErrs.ErrorLanguageNotFound)
		}

		return 0, fmt.Errorf("failed to get language id: %w", err)
	}

	course.LanguageID = languageID

	// check if course with same title already exists for this owner
	exists, err = s.repo.IsCourseExistsByOwnerAndTitle(ctx, course.OwnerID, course.Title)
	if err != nil {
		return 0, fmt.Errorf("failed to check course existence: %w", err)
	}

	if exists {
		return 0, serviceErrs.ErrorCourseExists
	}

	// create transactionally state machine item + course in database
	courseID, err := s.repo.CreateCourse(ctx, course)
	if err != nil {
		return 0, fmt.Errorf("failed to create course: %w", err)
	}

	return courseID, nil
}
