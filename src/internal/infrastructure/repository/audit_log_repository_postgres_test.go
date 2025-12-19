package repository_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/config"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/infrastructure/database"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/infrastructure/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, *entity.User) {
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Skipf("skipping audit log repo tests: %v", err)
	}

	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		t.Skipf("skipping audit log repo tests: %v", err)
	}

	gormDB := db.GetDB()

	err = database.RunMigrations(gormDB)
	assert.NoError(t, err)

	t.Cleanup(func() {
		database.DropAllTables(gormDB)
		db.Close()
	})

	user, err := entity.NewUser("test@example.com", "password123", "Test", "User", entity.UserRoleStudent)
	assert.NoError(t, err)
	err = gormDB.Create(user).Error
	assert.NoError(t, err)

	return gormDB, user
}

func TestNewPostgresAuditLogRepository(t *testing.T) {
	db, _ := setupTestDB(t)
	repo := repository.NewPostgresAuditLogRepository(db)
	assert.NotNil(t, repo)
}

func TestPostgresAuditLogRepository_CreateAndGetByID(t *testing.T) {
	db, user := setupTestDB(t)
	repo := repository.NewPostgresAuditLogRepository(db)

	auditLog, err := entity.NewAuditLog(entity.AuditActionCourseCreated, uuid.New(), "course", "{}", "{}", &user.ID)
	assert.NoError(t, err)

	err = repo.Create(context.Background(), auditLog)
	assert.NoError(t, err)

	found, err := repo.GetByID(context.Background(), auditLog.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, auditLog.ID, found.ID)
}

func TestPostgresAuditLogRepository_GetByID_NotFound(t *testing.T) {
	db, _ := setupTestDB(t)
	repo := repository.NewPostgresAuditLogRepository(db)

	_, err := repo.GetByID(context.Background(), uuid.New())
	assert.Error(t, err)
	assert.Equal(t, entity.ErrNotFound, err)
}

func TestPostgresAuditLogRepository_List(t *testing.T) {
	db, user := setupTestDB(t)
	repo := repository.NewPostgresAuditLogRepository(db)

	auditLog1, err := entity.NewAuditLog(entity.AuditActionCourseCreated, uuid.New(), "course", "{}", "{}", &user.ID)
	assert.NoError(t, err)
	err = repo.Create(context.Background(), auditLog1)
	assert.NoError(t, err)

	auditLog2, err := entity.NewAuditLog(entity.AuditActionCourseCreated, uuid.New(), "course", "{}", "{}", &user.ID)
	assert.NoError(t, err)
	err = repo.Create(context.Background(), auditLog2)
	assert.NoError(t, err)

	logs, total, err := repo.List(context.Background(), 10, 0)
	assert.NoError(t, err)
	assert.NotNil(t, logs)
	assert.GreaterOrEqual(t, len(logs), 2)
	assert.Equal(t, int64(len(logs)), total)
}

func TestPostgresAuditLogRepository_List_WithPagination(t *testing.T) {
	db, user := setupTestDB(t)
	repo := repository.NewPostgresAuditLogRepository(db)

	// Create 3 audit logs
	for i := 0; i < 3; i++ {
		auditLog, err := entity.NewAuditLog(entity.AuditActionCourseCreated, uuid.New(), "course", "{}", "{}", &user.ID)
		assert.NoError(t, err)
		err = repo.Create(context.Background(), auditLog)
		assert.NoError(t, err)
	}

	// Get first page
	logs, total, err := repo.List(context.Background(), 2, 0)
	assert.NoError(t, err)
	assert.NotNil(t, logs)
	assert.Len(t, logs, 2)
	assert.Equal(t, int64(3), total)

	// Get second page
	logs, total, err = repo.List(context.Background(), 2, 2)
	assert.NoError(t, err)
	assert.NotNil(t, logs)
	assert.GreaterOrEqual(t, len(logs), 1)
	assert.Equal(t, int64(3), total)
}

func TestPostgresAuditLogRepository_List_Empty(t *testing.T) {
	db, _ := setupTestDB(t)
	repo := repository.NewPostgresAuditLogRepository(db)

	// Clear any existing logs from setup
	err := db.Exec("DELETE FROM audit_logs").Error
	assert.NoError(t, err)

	logs, total, err := repo.List(context.Background(), 10, 0)
	assert.NoError(t, err)
	assert.NotNil(t, logs)
	assert.Len(t, logs, 0)
	assert.Equal(t, int64(0), total)
}
