package repository

import (
	"context"

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
	GetByStudentAndCourse(ctx context.Context, studentID, courseID string) (*entity.Enrollment, error)
	Update(ctx context.Context, enrollment *entity.Enrollment) error
	Delete(ctx context.Context, studentID, courseID string) error
	GetByStudent(ctx context.Context, studentID string, filter *EnrollmentFilter) ([]*entity.Enrollment, error)
	GetByCourse(ctx context.Context, courseID string, filter *EnrollmentFilter) ([]*entity.Enrollment, error)
}
