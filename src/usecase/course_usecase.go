package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/repository"
)

type CourseUseCase struct {
	courseRepo repository.CourseRepository
	userRepo   repository.UserRepository
}

func NewCourseUseCase(courseRepo repository.CourseRepository, userRepo repository.UserRepository) *CourseUseCase {
	return &CourseUseCase{
		courseRepo: courseRepo,
		userRepo:   userRepo,
	}
}

func (uc *CourseUseCase) CreateCourse(ctx context.Context, title, description, instructorID string, difficultyLevel entity.DifficultyLevel) (*entity.Course, error) {
	instructor, err := uc.userRepo.GetByID(ctx, instructorID)
	if err != nil {
		return nil, err
	}

	if !instructor.IsInstructor() {
		return nil, entity.ErrInsufficientPermissions
	}

	// Default to beginner if not provided
	if difficultyLevel == "" {
		difficultyLevel = entity.DifficultyLevelBeginner
	}

	course := &entity.Course{
		ID:              uuid.New().String(),
		Title:           title,
		Description:     description,
		InstructorID:    instructorID,
		DifficultyLevel: difficultyLevel,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := uc.courseRepo.Create(ctx, course); err != nil {
		return nil, err
	}

	return course, nil
}

func (uc *CourseUseCase) GetCourseByID(ctx context.Context, id string) (*entity.Course, error) {
	return uc.courseRepo.GetByID(ctx, id)
}

func (uc *CourseUseCase) ListCourses(ctx context.Context, filter *repository.CourseFilter) ([]*entity.Course, error) {
	return uc.courseRepo.List(ctx, filter)
}

func (uc *CourseUseCase) GetCoursesByInstructor(ctx context.Context, instructorID string) ([]*entity.Course, error) {
	return uc.courseRepo.GetByInstructor(ctx, instructorID)
}

func (uc *CourseUseCase) UpdateCourse(ctx context.Context, course *entity.Course, userID string) error {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if !user.IsAdmin() && course.InstructorID != userID {
		return entity.ErrInsufficientPermissions
	}
	course.UpdatedAt = time.Now()
	return uc.courseRepo.Update(ctx, course)
}

func (uc *CourseUseCase) DeleteCourse(ctx context.Context, courseID, userID string) error {
	course, err := uc.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return err
	}
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if !user.IsAdmin() && course.InstructorID != userID {
		return entity.ErrInsufficientPermissions
	}
	return uc.courseRepo.Delete(ctx, courseID)
}
