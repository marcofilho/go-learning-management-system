package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/repository"
)

type PostgresModuleRepository struct {
	db *gorm.DB
}

func NewPostgresModuleRepository(db *gorm.DB) repository.ModuleRepository {
	return &PostgresModuleRepository{db: db}
}

func (r *PostgresModuleRepository) Create(ctx context.Context, module *entity.Module) error {
	return r.db.WithContext(ctx).Create(module).Error
}

func (r *PostgresModuleRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Module, error) {
	var module entity.Module
	if err := r.db.WithContext(ctx).First(&module, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}
	return &module, nil
}

func (r *PostgresModuleRepository) GetByCourse(ctx context.Context, courseID uuid.UUID) ([]*entity.Module, error) {
	var modules []*entity.Module
	if err := r.db.WithContext(ctx).
		Where("course_id = ?", courseID).
		Order("order_index ASC, created_at ASC").
		Find(&modules).Error; err != nil {
		return nil, err
	}
	return modules, nil
}

func (r *PostgresModuleRepository) Update(ctx context.Context, module *entity.Module) error {
	return r.db.WithContext(ctx).Save(module).Error
}

func (r *PostgresModuleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Module{}, "id = ?", id).Error
}
