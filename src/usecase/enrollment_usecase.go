package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/repository"
)

type EnrollmentUseCase struct {
	enrollmentRepo repository.EnrollmentRepository
	courseRepo     repository.CourseRepository
	userRepo       repository.UserRepository
	auditLogRepo   repository.AuditLogRepository
}

func NewEnrollmentUseCase(enrollmentRepo repository.EnrollmentRepository, courseRepo repository.CourseRepository, userRepo repository.UserRepository, auditLogRepo repository.AuditLogRepository) *EnrollmentUseCase {
	return &EnrollmentUseCase{
		enrollmentRepo: enrollmentRepo,
		courseRepo:     courseRepo,
		userRepo:       userRepo,
		auditLogRepo:   auditLogRepo,
	}
}

func (uc *EnrollmentUseCase) EnrollStudent(ctx context.Context, studentID, courseID, requestorID uuid.UUID) (*entity.Enrollment, error) {
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

	// Create audit log
	payloadAfter := fmt.Sprintf(`{"student_id":"%s","course_id":"%s","status":"%s"}`, studentID, courseID, enrollment.Status)
	auditLog, _ := entity.NewAuditLog(entity.AuditActionEnrollmentCreated, courseID, "enrollment", "", payloadAfter, &requestorID)
	uc.auditLogRepo.Create(ctx, auditLog)

	return enrollment, nil
}

func (uc *EnrollmentUseCase) GetStudentCourses(ctx context.Context, studentID uuid.UUID, filter *repository.EnrollmentFilter) ([]*entity.Enrollment, int64, error) {
	return uc.enrollmentRepo.GetByStudent(ctx, studentID, filter)
}

func (uc *EnrollmentUseCase) GetCourseStudents(ctx context.Context, courseID uuid.UUID, filter *repository.EnrollmentFilter) ([]*entity.Enrollment, int64, error) {
	return uc.enrollmentRepo.GetByCourse(ctx, courseID, filter)
}

func (uc *EnrollmentUseCase) GetEnrollment(ctx context.Context, studentID, courseID uuid.UUID) (*entity.Enrollment, error) {
	return uc.enrollmentRepo.GetByStudentAndCourse(ctx, studentID, courseID)
}

func (uc *EnrollmentUseCase) UpdateEnrollmentStatus(ctx context.Context, studentID, courseID uuid.UUID, status entity.EnrollmentStatus, userID uuid.UUID) (*entity.Enrollment, error) {
	enrollment, err := uc.enrollmentRepo.GetByStudentAndCourse(ctx, studentID, courseID)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if !enrollment.CanBeModifiedBy(user) {
		return nil, entity.ErrUnauthorized
	}

	oldStatus := enrollment.Status
	enrollment.Status = status
	enrollment.UpdatedAt = time.Now()

	if err := uc.enrollmentRepo.Update(ctx, enrollment); err != nil {
		return nil, err
	}

	// Create audit log
	payloadBefore := fmt.Sprintf(`{"status":"%s"}`, oldStatus)
	payloadAfter := fmt.Sprintf(`{"status":"%s"}`, status)
	auditLog, _ := entity.NewAuditLog(entity.AuditActionEnrollmentUpdated, courseID, "enrollment", payloadBefore, payloadAfter, &userID)
	uc.auditLogRepo.Create(ctx, auditLog)

	return enrollment, nil
}

func (uc *EnrollmentUseCase) DropEnrollment(ctx context.Context, studentID, courseID, userID uuid.UUID) (*entity.Enrollment, error) {
	return uc.UpdateEnrollmentStatus(ctx, studentID, courseID, entity.EnrollmentStatusDropped, userID)
}

func (uc *EnrollmentUseCase) CompleteEnrollment(ctx context.Context, studentID, courseID, userID uuid.UUID) (*entity.Enrollment, error) {
	return uc.UpdateEnrollmentStatus(ctx, studentID, courseID, entity.EnrollmentStatusCompleted, userID)
}
