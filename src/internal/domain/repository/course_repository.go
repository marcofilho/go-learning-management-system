package repository

import (
	"context"

	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
)

type CourseFilter struct {
	InstructorID    *string
	DifficultyLevel *entity.DifficultyLevel
	ActiveOnly      bool
	Limit           int
	Offset          int
}

type CourseRepository interface {
	Create(ctx context.Context, course *entity.Course) error
	GetByID(ctx context.Context, id string) (*entity.Course, error)
	Update(ctx context.Context, course *entity.Course) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter *CourseFilter) ([]*entity.Course, error)
	GetByInstructor(ctx context.Context, instructorID string) ([]*entity.Course, error)
}
