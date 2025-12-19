package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
)

type AuditLogRepository interface {
	Create(ctx context.Context, auditLog *entity.AuditLog) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.AuditLog, error)
	List(ctx context.Context, limit, offset int) ([]*entity.AuditLog, int64, error)
}
