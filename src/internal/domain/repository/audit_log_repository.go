package repository

import (
	"context"

	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
)

type AuditLogRepository interface {
	Create(ctx context.Context, auditLog *entity.AuditLog) error
	GetByID(ctx context.Context, id string) (*entity.AuditLog, error)
	List(ctx context.Context, limit, offset int) ([]*entity.AuditLog, error)
}
