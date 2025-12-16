package repository

import (
	"context"

	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
)

type ModuleRepository interface {
	Create(ctx context.Context, module *entity.Module) error
	GetByID(ctx context.Context, id string) (*entity.Module, error)
	GetByCourse(ctx context.Context, courseID string) ([]*entity.Module, error)
	Update(ctx context.Context, module *entity.Module) error
	Delete(ctx context.Context, id string) error
}
