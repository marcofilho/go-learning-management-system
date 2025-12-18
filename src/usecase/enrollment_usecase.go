package usecase

import (
	"context"
	"time"

	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/repository"
)

type EnrollmentUseCase struct {
	enrollmentRepo repository.EnrollmentRepository
	courseRepo     repository.CourseRepository
	userRepo       repository.UserRepository
}

func NewEnrollmentUseCase(enrollmentRepo repository.EnrollmentRepository, courseRepo repository.CourseRepository, userRepo repository.UserRepository) *EnrollmentUseCase {
	return &EnrollmentUseCase{
		enrollmentRepo: enrollmentRepo,
		courseRepo:     courseRepo,
		userRepo:       userRepo,
	}
}

func (uc *EnrollmentUseCase) EnrollStudent(ctx context.Context, studentID, courseID, requestorID string) (*entity.Enrollment, error) {
	if _, err := uc.courseRepo.GetByID(ctx, courseID); err != nil {
		return nil, err
	}

	student, err := uc.userRepo.GetByID(ctx, studentID)
	if err != nil {
		return nil, err
	}

	requestor, err := uc.userRepo.GetByID(ctx, requestorID)
	if err != nil {
		return nil, err
	}

	if !requestor.IsAdmin() && studentID != requestorID {
		return nil, entity.ErrUnauthorized
	}

	existing, _ := uc.enrollmentRepo.GetByStudentAndCourse(ctx, studentID, courseID)
	if existing != nil {
		return nil, entity.ErrAlreadyEnrolled
	}

	enrollment, err := entity.NewEnrollment(student, courseID)
	if err != nil {
		return nil, err
	}

	if err := uc.enrollmentRepo.Create(ctx, enrollment); err != nil {
		return nil, err
	}

	return enrollment, nil
}

func (uc *EnrollmentUseCase) GetStudentCourses(ctx context.Context, studentID string, filter *repository.EnrollmentFilter) ([]*entity.Enrollment, error) {
	return uc.enrollmentRepo.GetByStudent(ctx, studentID, filter)
}

func (uc *EnrollmentUseCase) GetCourseStudents(ctx context.Context, courseID string, filter *repository.EnrollmentFilter) ([]*entity.Enrollment, error) {
	return uc.enrollmentRepo.GetByCourse(ctx, courseID, filter)
}

func (uc *EnrollmentUseCase) UpdateEnrollmentStatus(ctx context.Context, studentID, courseID string, status entity.EnrollmentStatus, userID string) error {
	enrollment, err := uc.enrollmentRepo.GetByStudentAndCourse(ctx, studentID, courseID)
	if err != nil {
		return err
	}

	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if !enrollment.CanBeModifiedBy(user) {
		return entity.ErrUnauthorized
	}

	enrollment.Status = status
	enrollment.UpdatedAt = time.Now()

	return uc.enrollmentRepo.Update(ctx, enrollment)
}

func (uc *EnrollmentUseCase) DropEnrollment(ctx context.Context, studentID, courseID, userID string) error {
	return uc.UpdateEnrollmentStatus(ctx, studentID, courseID, entity.EnrollmentStatusDropped, userID)
}

func (uc *EnrollmentUseCase) CompleteEnrollment(ctx context.Context, studentID, courseID, userID string) error {
	return uc.UpdateEnrollmentStatus(ctx, studentID, courseID, entity.EnrollmentStatusCompleted, userID)
}
