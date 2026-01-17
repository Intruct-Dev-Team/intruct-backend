package course

import (
	"context"
	"errors"
	"fmt"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/entities"
	serviceErrs "github.com/Intruct-Dev-Team/intruct-backend/internal/errors"
	"github.com/Intruct-Dev-Team/intruct-backend/pkg/workerpool"
)

func (s *CourseService) GetCourseList(ctx context.Context, userID int, inMine bool) ([]*entities.Course, error) {
	ids, err := s.repo.GetOwnCourseIDs(ctx, userID)
	if err != nil && !errors.Is(err, serviceErrs.ErrorSelectEmpty) {
		return nil, fmt.Errorf("failed to get own course IDs: %w", err)
	}

	courseMap := make(map[int]*entities.Course)
	for _, id := range ids {
		courseMap[id] = nil
	}

	if inMine {
		// only my courses + study courses
		ids, err = s.repo.GetCourseIDsFromUserProgress(ctx, userID)
		if err != nil && !errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return nil, fmt.Errorf("failed to get study course IDs: %w", err)
		}
		for _, id := range ids {
			courseMap[id] = nil
		}

	} else {
		// all courses
		ids, err = s.repo.GetPublicCourseIDs(ctx)
		if err != nil && !errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return nil, fmt.Errorf("failed to get public course IDs: %w", err)
		}
		for _, id := range ids {
			courseMap[id] = nil
		}
	}

	if len(courseMap) == 0 {
		return []*entities.Course{}, nil
	}

	// get all needed courses
	wp := workerpool.NewWorkerPool[entities.Course](5)
	err = wp.FillMap(ctx, courseMap, s.repo.GetCourseByID)
	if err != nil {
		return nil, fmt.Errorf("failed to get courses by id: %w", err)
	}

	// get user's course progressions
	progressions, err := s.repo.GetUsersCourseProgressions(ctx, userID)
	if err != nil && !errors.Is(err, serviceErrs.ErrorSelectEmpty) {
		return nil, fmt.Errorf("failed to get user progressions: %w", err)
	}

	for _, p := range progressions {
		course, ok := courseMap[p.CourseID]
		if ok {
			course.CourseProgression = p
		}
	}

	courses := make([]*entities.Course, 0, len(courseMap))
	for _, course := range courseMap {
		courses = append(courses, course)
	}

	return courses, nil
}
