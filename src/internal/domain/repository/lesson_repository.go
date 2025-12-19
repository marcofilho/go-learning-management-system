package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
)

type LessonRepository interface {
	CreateLesson(ctx context.Context, lesson *entity.Lesson) error
	GetLessonByID(ctx context.Context, id uuid.UUID) (*entity.Lesson, error)
	DeleteLesson(ctx context.Context, lessonID uuid.UUID) error

	CreateVersion(ctx context.Context, lesson *entity.LessonVersion) error
	GetVersionByID(ctx context.Context, id uuid.UUID) (*entity.LessonVersion, error)
	GetLatestByModule(ctx context.Context, moduleID uuid.UUID, limit, offset int) ([]*entity.LessonVersion, int64, error)
	GetAllVersions(ctx context.Context, lessonID uuid.UUID, limit, offset int) ([]*entity.LessonVersion, int64, error)
	GetNextVersionNumber(ctx context.Context, lessonID uuid.UUID) (int, error)
	DeleteVersionsByLesson(ctx context.Context, lessonID uuid.UUID) error
}
