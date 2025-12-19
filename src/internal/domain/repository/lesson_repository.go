package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
)

type LessonRepository interface {
	Create(ctx context.Context, lesson *entity.LessonVersion) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.LessonVersion, error)
	GetByModule(ctx context.Context, moduleID uuid.UUID) ([]*entity.LessonVersion, error)
	GetLatestByModule(ctx context.Context, moduleID uuid.UUID) ([]*entity.LessonVersion, error)
	GetAllVersions(ctx context.Context, lessonID uuid.UUID) ([]*entity.LessonVersion, error)
	GetNextVersionNumber(ctx context.Context, moduleID uuid.UUID) (int, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
