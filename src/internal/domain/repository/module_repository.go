package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
)

type ModuleRepository interface {
	Create(ctx context.Context, module *entity.Module) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Module, error)
	GetByCourse(ctx context.Context, courseID uuid.UUID) ([]*entity.Module, error)
	Update(ctx context.Context, module *entity.Module) error
	Delete(ctx context.Context, id uuid.UUID) error
}
