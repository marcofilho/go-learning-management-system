package repository

import (
	"context"

	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
)

type LessonRepository interface {
	Create(ctx context.Context, lesson *entity.LessonVersion) error
	GetByID(ctx context.Context, id string) (*entity.LessonVersion, error)
	GetByModule(ctx context.Context, moduleID string) ([]*entity.LessonVersion, error)
	GetLatestByModule(ctx context.Context, moduleID string) ([]*entity.LessonVersion, error)
	GetAllVersions(ctx context.Context, lessonID string) ([]*entity.LessonVersion, error)
	GetNextVersionNumber(ctx context.Context, moduleID string) (int, error)
	Delete(ctx context.Context, id string) error
}
