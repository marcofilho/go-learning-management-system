package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
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
	return r.db.WithContext(ctx).Exec(
		`INSERT INTO audit_logs (id, action, resource_id, resource_type, user_id, payload_before, payload_after, created_at) 
		 VALUES (?, ?, ?, ?, ?, ?::jsonb, ?::jsonb, ?)`,
		auditLog.ID,
		auditLog.Action,
		auditLog.EntityID.String(),
		auditLog.EntityType,
		auditLog.UserID,
		auditLog.OldValues,
		auditLog.NewValues,
		auditLog.CreatedAt,
	).Error
}

func (r *PostgresAuditLogRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.AuditLog, error) {
	var result struct {
		ID           uuid.UUID
		Action       string
		ResourceID   string
		ResourceType string
		UserID       sql.NullString
		OldValues    sql.NullString
		NewValues    sql.NullString
		CreatedAt    time.Time
	}

	if err := r.db.WithContext(ctx).Raw(
		`SELECT id, action, resource_id, resource_type, user_id, payload_before, payload_after, created_at 
		 FROM audit_logs WHERE id = ?`, id,
	).Scan(&result).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}

	entityID, err := uuid.Parse(result.ResourceID)
	if err != nil {
		return nil, err
	}

	auditLog := &entity.AuditLog{
		ID:         result.ID,
		Action:     entity.AuditAction(result.Action),
		EntityID:   entityID,
		EntityType: result.ResourceType,
		CreatedAt:  result.CreatedAt,
	}

	if result.UserID.Valid {
		userID, err := uuid.Parse(result.UserID.String)
		if err == nil {
			auditLog.UserID = &userID
		}
	}

	if result.OldValues.Valid {
		auditLog.OldValues = result.OldValues.String
	}

	if result.NewValues.Valid {
		auditLog.NewValues = result.NewValues.String
	}

	return auditLog, nil
}

func (r *PostgresAuditLogRepository) List(ctx context.Context, limit, offset int) ([]*entity.AuditLog, int64, error) {
	var results []struct {
		ID           uuid.UUID
		Action       string
		ResourceID   string
		ResourceType string
		UserID       sql.NullString
		OldValues    sql.NullString
		NewValues    sql.NullString
		CreatedAt    time.Time
	}

	var total int64
	if err := r.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM audit_logs`).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Raw(
		`SELECT id, action, resource_id, resource_type, user_id, payload_before, payload_after, created_at 
		 FROM audit_logs ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		limit, offset,
	).Scan(&results).Error; err != nil {
		return nil, 0, err
	}

	auditLogs := make([]*entity.AuditLog, len(results))
	for i, result := range results {
		entityID, err := uuid.Parse(result.ResourceID)
		if err != nil {
			continue
		}

		auditLog := &entity.AuditLog{
			ID:         result.ID,
			Action:     entity.AuditAction(result.Action),
			EntityID:   entityID,
			EntityType: result.ResourceType,
			CreatedAt:  result.CreatedAt,
		}

		if result.UserID.Valid {
			userID, err := uuid.Parse(result.UserID.String)
			if err == nil {
				auditLog.UserID = &userID
			}
		}

		if result.OldValues.Valid {
			auditLog.OldValues = result.OldValues.String
		}

		if result.NewValues.Valid {
			auditLog.NewValues = result.NewValues.String
		}

		auditLogs[i] = auditLog
	}

	return auditLogs, total, nil
}
