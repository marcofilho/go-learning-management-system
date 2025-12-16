package repository

import (
	"context"

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

func (r *PostgresLessonRepository) Create(ctx context.Context, lesson *entity.LessonVersion) error {
	return r.db.WithContext(ctx).Create(lesson).Error
}

func (r *PostgresLessonRepository) GetByID(ctx context.Context, id string) (*entity.LessonVersion, error) {
	var lesson entity.LessonVersion
	if err := r.db.WithContext(ctx).First(&lesson, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}
	return &lesson, nil
}

func (r *PostgresLessonRepository) GetByModule(ctx context.Context, moduleID string) ([]*entity.LessonVersion, error) {
	var lessons []*entity.LessonVersion
	if err := r.db.WithContext(ctx).
		Where("module_id = ?", moduleID).
		Order("version_number DESC, created_at DESC").
		Find(&lessons).Error; err != nil {
		return nil, err
	}
	return lessons, nil
}

func (r *PostgresLessonRepository) GetLatestByModule(ctx context.Context, moduleID string) ([]*entity.LessonVersion, error) {
	var lessons []*entity.LessonVersion

	// Get the latest version for each lesson in the module
	subQuery := r.db.Model(&entity.LessonVersion{}).
		Select("MAX(version_number)").
		Where("module_id = ?", moduleID).
		Group("module_id")

	if err := r.db.WithContext(ctx).
		Where("module_id = ? AND version_number IN (?)", moduleID, subQuery).
		Order("created_at DESC").
		Find(&lessons).Error; err != nil {
		return nil, err
	}
	return lessons, nil
}

func (r *PostgresLessonRepository) GetAllVersions(ctx context.Context, lessonID string) ([]*entity.LessonVersion, error) {
	var lesson entity.LessonVersion
	if err := r.db.WithContext(ctx).First(&lesson, "id = ?", lessonID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}

	var versions []*entity.LessonVersion
	if err := r.db.WithContext(ctx).
		Where("module_id = ?", lesson.ModuleID).
		Order("version_number DESC").
		Find(&versions).Error; err != nil {
		return nil, err
	}
	return versions, nil
}

func (r *PostgresLessonRepository) GetNextVersionNumber(ctx context.Context, moduleID string) (int, error) {
	var maxVersion int
	err := r.db.WithContext(ctx).
		Model(&entity.LessonVersion{}).
		Where("module_id = ?", moduleID).
		Select("COALESCE(MAX(version_number), 0)").
		Scan(&maxVersion).Error

	if err != nil {
		return 0, err
	}
	return maxVersion + 1, nil
}

func (r *PostgresLessonRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.LessonVersion{}, "id = ?", id).Error
}
