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
	assert.NoError(t, err)

	db, err := database.NewPostgresDB(&cfg.Database)
	assert.NoError(t, err)

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

func TestPostgresAuditLogRepository(t *testing.T) {
	db, user := setupTestDB(t)
	repo := repository.NewPostgresAuditLogRepository(db)

	t.Run("Create and GetByID", func(t *testing.T) {
		t.Parallel()
		auditLog, err := entity.NewAuditLog(entity.AuditActionCourseCreated, uuid.New().String(), "course", "{}", "{}", &user.ID)
		assert.NoError(t, err)

		err = repo.Create(context.Background(), auditLog)
		assert.NoError(t, err)

		found, err := repo.GetByID(context.Background(), auditLog.ID)
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, auditLog.ID, found.ID)
	})

	t.Run("GetByID not found", func(t *testing.T) {
		t.Parallel()
		_, err := repo.GetByID(context.Background(), uuid.New().String())
		assert.Error(t, err)
	})

	t.Run("List", func(t *testing.T) {
		t.Parallel()
		auditLog1, err := entity.NewAuditLog(entity.AuditActionCourseCreated, uuid.New().String(), "course", "{}", "{}", &user.ID)
		assert.NoError(t, err)
		err = repo.Create(context.Background(), auditLog1)
		assert.NoError(t, err)

		auditLog2, err := entity.NewAuditLog(entity.AuditActionCourseCreated, uuid.New().String(), "course", "{}", "{}", &user.ID)
		assert.NoError(t, err)
		err = repo.Create(context.Background(), auditLog2)
		assert.NoError(t, err)

		logs, err := repo.List(context.Background(), 10, 0)
		assert.NoError(t, err)
		assert.NotNil(t, logs)
		assert.GreaterOrEqual(t, len(logs), 2)
	})

	t.Run("List with pagination", func(t *testing.T) {
		t.Parallel()
		// Create 3 audit logs
		for i := 0; i < 3; i++ {
			auditLog, err := entity.NewAuditLog(entity.AuditActionCourseCreated, uuid.New().String(), "course", "{}", "{}", &user.ID)
			assert.NoError(t, err)
			err = repo.Create(context.Background(), auditLog)
			assert.NoError(t, err)
		}

		// Get first page
		logs, err := repo.List(context.Background(), 2, 0)
		assert.NoError(t, err)
		assert.NotNil(t, logs)
		assert.Len(t, logs, 2)

		// Get second page
		logs, err = repo.List(context.Background(), 2, 2)
		assert.NoError(t, err)
		assert.NotNil(t, logs)
		assert.GreaterOrEqual(t, len(logs), 1)
	})
}
