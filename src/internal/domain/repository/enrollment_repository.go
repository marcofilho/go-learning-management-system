package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
)

type EnrollmentFilter struct {
	Status   *entity.EnrollmentStatus
	DateFrom *string
	DateTo   *string
	Limit    int
	Offset   int
}

type EnrollmentRepository interface {
	Create(ctx context.Context, enrollment *entity.Enrollment) error
	GetByStudentAndCourse(ctx context.Context, studentID, courseID uuid.UUID) (*entity.Enrollment, error)
	Update(ctx context.Context, enrollment *entity.Enrollment) error
	Delete(ctx context.Context, studentID, courseID uuid.UUID) error
	GetByStudent(ctx context.Context, studentID uuid.UUID, filter *EnrollmentFilter) ([]*entity.Enrollment, int64, error)
	GetByCourse(ctx context.Context, courseID uuid.UUID, filter *EnrollmentFilter) ([]*entity.Enrollment, int64, error)
}
