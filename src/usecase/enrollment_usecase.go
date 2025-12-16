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
	// Verify course exists
	if _, err := uc.courseRepo.GetByID(ctx, courseID); err != nil {
		return nil, err
	}

	// Verify student exists and is active
	student, err := uc.userRepo.GetByID(ctx, studentID)
	if err != nil {
		return nil, err
	}

	if !student.IsActive {
		return nil, entity.ErrUnauthorized
	}

	// Students must have role student
	if student.Role != entity.UserRoleStudent {
		return nil, entity.ErrInvalidInput
	}

	// Check authorization: students can enroll themselves, admins can enroll anyone
	requestor, err := uc.userRepo.GetByID(ctx, requestorID)
	if err != nil {
		return nil, err
	}

	if !requestor.IsAdmin() && studentID != requestorID {
		return nil, entity.ErrUnauthorized
	}

	// Check if already enrolled
	existing, _ := uc.enrollmentRepo.GetByStudentAndCourse(ctx, studentID, courseID)
	if existing != nil {
		return nil, entity.ErrAlreadyEnrolled
	}

	enrollment := &entity.Enrollment{
		StudentID:      studentID,
		CourseID:       courseID,
		EnrollmentDate: time.Now(),
		Status:         entity.EnrollmentStatusActive,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
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

	// Students can update their own enrollment, admins can update any
	if !user.IsAdmin() && studentID != userID {
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
