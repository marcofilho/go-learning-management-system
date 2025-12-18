package repository

import (
	"context"

	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/repository"
	"gorm.io/gorm"
)

type PostgresAuditLogRepository struct {
	db *gorm.DB
}

func NewPostgresAuditLogRepository(db *gorm.DB) repository.AuditLogRepository {
	return &PostgresAuditLogRepository{db: db}
}

func (r *PostgresAuditLogRepository) Create(ctx context.Context, auditLog *entity.AuditLog) error {
	return r.db.WithContext(ctx).Create(auditLog).Error
}

func (r *PostgresAuditLogRepository) GetByID(ctx context.Context, id string) (*entity.AuditLog, error) {
	var auditLog entity.AuditLog
	if err := r.db.WithContext(ctx).First(&auditLog, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}
	return &auditLog, nil
}

func (r *PostgresAuditLogRepository) List(ctx context.Context, limit, offset int) ([]*entity.AuditLog, error) {
	var auditLogs []*entity.AuditLog
	if err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&auditLogs).Error; err != nil {
		return nil, err
	}
	return auditLogs, nil
}
