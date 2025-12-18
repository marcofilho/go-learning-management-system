package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/repository"
)

type CourseUseCase struct {
	courseRepo   repository.CourseRepository
	userRepo     repository.UserRepository
	auditLogRepo repository.AuditLogRepository
}

func NewCourseUseCase(courseRepo repository.CourseRepository, userRepo repository.UserRepository, auditLogRepo repository.AuditLogRepository) *CourseUseCase {
	return &CourseUseCase{
		courseRepo:   courseRepo,
		userRepo:     userRepo,
		auditLogRepo: auditLogRepo,
	}
}

func (uc *CourseUseCase) CreateCourse(ctx context.Context, title, description, instructorID string, difficultyLevel entity.DifficultyLevel) (*entity.Course, error) {
	instructor, err := uc.userRepo.GetByID(ctx, instructorID)
	if err != nil {
		return nil, entity.ErrNotFound
	}

	course, err := entity.NewCourse(title, description, instructorID, difficultyLevel, instructor)
	if err := uc.courseRepo.Create(ctx, course); err != nil {
		return nil, err
	}

	// Create audit log
	payloadAfter := fmt.Sprintf(`{"id":"%s","title":"%s","instructor_id":"%s","difficulty_level":"%s"}`, course.ID, course.Title, course.InstructorID, course.DifficultyLevel)
	auditLog, _ := entity.NewAuditLog(entity.AuditActionCourseCreated, course.ID, "course", "", payloadAfter, &instructorID)
	uc.auditLogRepo.Create(ctx, auditLog)

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
	if err := course.Validate(); err != nil {
		return err
	}

	// Get existing course for audit log
	existingCourse, err := uc.courseRepo.GetByID(ctx, course.ID)
	if err != nil {
		return err
	}

	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if !course.CanBeModifiedBy(user) {
		return entity.ErrInsufficientPermissions
	}

	course.UpdatedAt = time.Now()
	if err := uc.courseRepo.Update(ctx, course); err != nil {
		return err
	}

	// Create audit log
	payloadBefore := fmt.Sprintf(`{"title":"%s","description":"%s","difficulty_level":"%s"}`, existingCourse.Title, existingCourse.Description, existingCourse.DifficultyLevel)
	payloadAfter := fmt.Sprintf(`{"title":"%s","description":"%s","difficulty_level":"%s"}`, course.Title, course.Description, course.DifficultyLevel)
	auditLog, _ := entity.NewAuditLog(entity.AuditActionCourseUpdated, course.ID, "course", payloadBefore, payloadAfter, &userID)
	uc.auditLogRepo.Create(ctx, auditLog)

	return nil
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

	if !course.CanBeModifiedBy(user) {
		return entity.ErrInsufficientPermissions
	}

	if err := uc.courseRepo.Delete(ctx, courseID); err != nil {
		return err
	}

	// Create audit log
	payloadBefore := fmt.Sprintf(`{"id":"%s","title":"%s","instructor_id":"%s"}`, course.ID, course.Title, course.InstructorID)
	auditLog, _ := entity.NewAuditLog(entity.AuditActionCourseDeleted, courseID, "course", payloadBefore, "", &userID)
	uc.auditLogRepo.Create(ctx, auditLog)

	return nil
}
