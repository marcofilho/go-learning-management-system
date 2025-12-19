package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/repository"
	"gorm.io/gorm"
)

type PostgresCourseRepository struct {
	db *gorm.DB
}

func NewPostgresCourseRepository(db *gorm.DB) repository.CourseRepository {
	return &PostgresCourseRepository{db: db}
}

func (r *PostgresCourseRepository) Create(ctx context.Context, course *entity.Course) error {
	return r.db.WithContext(ctx).Create(course).Error
}

func (r *PostgresCourseRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Course, error) {
	var course entity.Course
	if err := r.db.WithContext(ctx).First(&course, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}
	return &course, nil
}

func (r *PostgresCourseRepository) Update(ctx context.Context, course *entity.Course) error {
	return r.db.WithContext(ctx).Save(course).Error
}

func (r *PostgresCourseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Course{}, "id = ?", id).Error
}

func (r *PostgresCourseRepository) List(ctx context.Context, filter *repository.CourseFilter) ([]*entity.Course, error) {
	var courses []*entity.Course
	query := r.db.WithContext(ctx)

	// Apply filters
	if filter.InstructorID != nil {
		query = query.Where("instructor_id = ?", *filter.InstructorID)
	}

	if filter.DifficultyLevel != nil {
		query = query.Where("difficulty_level = ?", *filter.DifficultyLevel)
	}

	if filter.ActiveOnly {
		query = query.Where("deleted_at IS NULL")
	}

	// Apply pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	if err := query.Order("created_at DESC").Find(&courses).Error; err != nil {
		return nil, err
	}
	return courses, nil
}

func (r *PostgresCourseRepository) GetByInstructor(ctx context.Context, instructorID uuid.UUID) ([]*entity.Course, error) {
	var courses []*entity.Course
	if err := r.db.WithContext(ctx).
		Where("instructor_id = ?", instructorID).
		Order("created_at DESC").
		Find(&courses).Error; err != nil {
		return nil, err
	}
	return courses, nil
}
