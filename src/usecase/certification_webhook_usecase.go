package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/repository"
)

type CertificationWebhookUseCase struct {
	userRepo       repository.UserRepository
	courseRepo     repository.CourseRepository
	enrollmentRepo repository.EnrollmentRepository
	auditLogRepo   repository.AuditLogRepository
}

func NewCertificationWebhookUseCase(
	userRepo repository.UserRepository,
	courseRepo repository.CourseRepository,
	enrollmentRepo repository.EnrollmentRepository,
	auditLogRepo repository.AuditLogRepository,
) *CertificationWebhookUseCase {
	return &CertificationWebhookUseCase{
		userRepo:       userRepo,
		courseRepo:     courseRepo,
		enrollmentRepo: enrollmentRepo,
		auditLogRepo:   auditLogRepo,
	}
}

type CertificationWebhookPayload struct {
	StudentID           string    `json:"student_id"`
	CourseID            string    `json:"course_id"`
	CertificationStatus string    `json:"certification_status"`
	Score               int       `json:"score"`
	Timestamp           time.Time `json:"timestamp"`
}

func (uc *CertificationWebhookUseCase) ProcessCertification(ctx context.Context, payload CertificationWebhookPayload) error {
	// Validate student exists
	if _, err := uc.userRepo.GetByID(ctx, payload.StudentID); err != nil {
		return fmt.Errorf("%w: student not found", entity.ErrNotFound)
	}

	// Validate course exists
	if _, err := uc.courseRepo.GetByID(ctx, payload.CourseID); err != nil {
		return fmt.Errorf("%w: course not found", entity.ErrNotFound)
	}

	// Get enrollment
	enrollment, err := uc.enrollmentRepo.GetByStudentAndCourse(ctx, payload.StudentID, payload.CourseID)
	if err != nil {
		return fmt.Errorf("%w: enrollment not found", entity.ErrNotFound)
	}

	// Only update enrollments with status = active
	if enrollment.Status != entity.EnrollmentStatusActive {
		return fmt.Errorf("%w: enrollment status must be active", entity.ErrInvalidInput)
	}

	// Update enrollment based on certification status
	if payload.CertificationStatus == "passed" {
		enrollment.Status = entity.EnrollmentStatusCompleted
		enrollment.UpdatedAt = time.Now()

		if err := uc.enrollmentRepo.Update(ctx, enrollment); err != nil {
			return err
		}
	}

	// Create audit log entry with payload details
	payloadBeforeJSON := fmt.Sprintf(`{"status":"%s"}`, entity.EnrollmentStatusActive)
	payloadAfterJSON := fmt.Sprintf(
		`{"status":"%s","certification_status":"%s","score":%d}`,
		enrollment.Status,
		payload.CertificationStatus,
		payload.Score,
	)

	// Use composite key (StudentID-CourseID) as entity ID for audit log
	entityID := fmt.Sprintf("%s-%s", payload.StudentID, payload.CourseID)
	auditLog, err := entity.NewAuditLog(
		entity.AuditActionCertificationWebhook,
		entityID,
		"enrollment",
		payloadBeforeJSON,
		payloadAfterJSON,
		&payload.StudentID,
	)
	if err != nil {
		return err
	}

	if err := uc.auditLogRepo.Create(ctx, auditLog); err != nil {
		return err
	}

	return nil
}
