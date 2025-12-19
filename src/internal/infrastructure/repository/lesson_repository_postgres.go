package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/repository"
)

type PostgresLessonRepository struct {
	db *gorm.DB
}

func NewPostgresLessonRepository(db *gorm.DB) repository.LessonRepository {
	return &PostgresLessonRepository{db: db}
}

func (r *PostgresLessonRepository) CreateLesson(ctx context.Context, lesson *entity.Lesson) error {
	return r.db.WithContext(ctx).Create(lesson).Error
}

func (r *PostgresLessonRepository) GetLessonByID(ctx context.Context, id uuid.UUID) (*entity.Lesson, error) {
	var lesson entity.Lesson
	if err := r.db.WithContext(ctx).First(&lesson, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}
	return &lesson, nil
}

func (r *PostgresLessonRepository) DeleteLesson(ctx context.Context, lessonID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Lesson{}, "id = ?", lessonID).Error
}

func (r *PostgresLessonRepository) CreateVersion(ctx context.Context, lesson *entity.LessonVersion) error {
	return r.db.WithContext(ctx).Create(lesson).Error
}

func (r *PostgresLessonRepository) GetVersionByID(ctx context.Context, id uuid.UUID) (*entity.LessonVersion, error) {
	var lesson entity.LessonVersion
	if err := r.db.WithContext(ctx).First(&lesson, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}
	return &lesson, nil
}

func (r *PostgresLessonRepository) GetLatestByModule(ctx context.Context, moduleID uuid.UUID, limit, offset int) ([]*entity.LessonVersion, int64, error) {
	var lessons []*entity.LessonVersion
	query := r.db.WithContext(ctx).Model(&entity.LessonVersion{}).
		Where("module_id = ?", moduleID)

	// Latest version per lesson thread
	sub := r.db.Model(&entity.LessonVersion{}).
		Select("lesson_id, MAX(version_number) AS max_version").
		Where("module_id = ?", moduleID).
		Group("lesson_id")

	query = query.Joins("JOIN (?) lv ON lv.lesson_id = lesson_versions.lesson_id AND lv.max_version = lesson_versions.version_number", sub)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Order("created_at DESC").Find(&lessons).Error; err != nil {
		return nil, 0, err
	}

	return lessons, total, nil
}

func (r *PostgresLessonRepository) GetAllVersions(ctx context.Context, lessonID uuid.UUID, limit, offset int) ([]*entity.LessonVersion, int64, error) {
	var versions []*entity.LessonVersion
	query := r.db.WithContext(ctx).Model(&entity.LessonVersion{}).
		Where("lesson_id = ?", lessonID)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Order("version_number DESC").Find(&versions).Error; err != nil {
		return nil, 0, err
	}

	return versions, total, nil
}

func (r *PostgresLessonRepository) GetNextVersionNumber(ctx context.Context, lessonID uuid.UUID) (int, error) {
	var maxVersion int
	err := r.db.WithContext(ctx).
		Model(&entity.LessonVersion{}).
		Where("lesson_id = ?", lessonID).
		Select("COALESCE(MAX(version_number), 0)").
		Scan(&maxVersion).Error

	if err != nil {
		return 0, err
	}
	return maxVersion + 1, nil
}

func (r *PostgresLessonRepository) DeleteVersionsByLesson(ctx context.Context, lessonID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("lesson_id = ?", lessonID).Delete(&entity.LessonVersion{}).Error
}
