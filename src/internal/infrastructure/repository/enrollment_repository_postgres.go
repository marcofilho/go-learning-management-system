package repository

import (
	"context"
	"time"

	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/repository"
	"gorm.io/gorm"
)

type PostgresEnrollmentRepository struct {
	db *gorm.DB
}

func NewPostgresEnrollmentRepository(db *gorm.DB) repository.EnrollmentRepository {
	return &PostgresEnrollmentRepository{db: db}
}

func (r *PostgresEnrollmentRepository) Create(ctx context.Context, enrollment *entity.Enrollment) error {
	return r.db.WithContext(ctx).Create(enrollment).Error
}

func (r *PostgresEnrollmentRepository) GetByStudentAndCourse(ctx context.Context, studentID, courseID string) (*entity.Enrollment, error) {
	var enrollment entity.Enrollment
	if err := r.db.WithContext(ctx).
		Where("student_id = ? AND course_id = ?", studentID, courseID).
		First(&enrollment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}
	return &enrollment, nil
}

func (r *PostgresEnrollmentRepository) Update(ctx context.Context, enrollment *entity.Enrollment) error {
	return r.db.WithContext(ctx).Save(enrollment).Error
}

func (r *PostgresEnrollmentRepository) Delete(ctx context.Context, studentID, courseID string) error {
	return r.db.WithContext(ctx).
		Where("student_id = ? AND course_id = ?", studentID, courseID).
		Delete(&entity.Enrollment{}).Error
}

func (r *PostgresEnrollmentRepository) GetByStudent(ctx context.Context, studentID string, filter *repository.EnrollmentFilter) ([]*entity.Enrollment, error) {
	var enrollments []*entity.Enrollment
	query := r.db.WithContext(ctx).Where("student_id = ?", studentID)

	if filter != nil {
		query = r.applyFilters(query, filter)
	}

	if err := query.Order("enrollment_date DESC").Find(&enrollments).Error; err != nil {
		return nil, err
	}
	return enrollments, nil
}

func (r *PostgresEnrollmentRepository) GetByCourse(ctx context.Context, courseID string, filter *repository.EnrollmentFilter) ([]*entity.Enrollment, error) {
	var enrollments []*entity.Enrollment
	query := r.db.WithContext(ctx).Where("course_id = ?", courseID)

	if filter != nil {
		query = r.applyFilters(query, filter)
	}

	if err := query.Order("enrollment_date DESC").Find(&enrollments).Error; err != nil {
		return nil, err
	}
	return enrollments, nil
}

func (r *PostgresEnrollmentRepository) applyFilters(query *gorm.DB, filter *repository.EnrollmentFilter) *gorm.DB {
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	if filter.DateFrom != nil {
		if dateFrom, err := time.Parse("2006-01-02", *filter.DateFrom); err == nil {
			query = query.Where("enrollment_date >= ?", dateFrom)
		}
	}

	if filter.DateTo != nil {
		if dateTo, err := time.Parse("2006-01-02", *filter.DateTo); err == nil {
			query = query.Where("enrollment_date <= ?", dateTo)
		}
	}

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}

	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	return query
}
