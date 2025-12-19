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

func setupUserTestDB(t *testing.T) *gorm.DB {
	t.Helper()

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

	return gormDB
}

func TestNewPostgresUserRepository(t *testing.T) {
	db := setupUserTestDB(t)
	repo := repository.NewPostgresUserRepository(db)
	assert.NotNil(t, repo)
}

func TestPostgresUserRepository_CreateAndGet(t *testing.T) {
	db := setupUserTestDB(t)
	repo := repository.NewPostgresUserRepository(db)

	user, err := entity.NewUser("test@example.com", "password123", "Test", "User", entity.UserRoleStudent)
	assert.NoError(t, err)

	err = repo.Create(context.Background(), user)
	assert.NoError(t, err)

	t.Run("GetByID", func(t *testing.T) {
		found, err := repo.GetByID(context.Background(), user.ID)
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, user.ID, found.ID)
	})

	t.Run("GetByEmail", func(t *testing.T) {
		found, err := repo.GetByEmail(context.Background(), user.Email)
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, user.Email, found.Email)
	})
}

func TestPostgresUserRepository_Update(t *testing.T) {
	db := setupUserTestDB(t)
	repo := repository.NewPostgresUserRepository(db)

	user, err := entity.NewUser("test@example.com", "password123", "Test", "User", entity.UserRoleStudent)
	assert.NoError(t, err)
	err = repo.Create(context.Background(), user)
	assert.NoError(t, err)

	user.FirstName = "Updated"
	err = repo.Update(context.Background(), user)
	assert.NoError(t, err)

	found, err := repo.GetByID(context.Background(), user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated", found.FirstName)
}

func TestPostgresUserRepository_Delete(t *testing.T) {
	db := setupUserTestDB(t)
	repo := repository.NewPostgresUserRepository(db)

	user, err := entity.NewUser("test@example.com", "password123", "Test", "User", entity.UserRoleStudent)
	assert.NoError(t, err)
	err = repo.Create(context.Background(), user)
	assert.NoError(t, err)

	err = repo.Delete(context.Background(), user.ID)
	assert.NoError(t, err)

	_, err = repo.GetByID(context.Background(), user.ID)
	assert.Error(t, err)
	assert.Equal(t, entity.ErrNotFound, err)
}

func TestPostgresUserRepository_List(t *testing.T) {
	db := setupUserTestDB(t)
	repo := repository.NewPostgresUserRepository(db)

	user1, _ := entity.NewUser("user1@example.com", "pass", "First", "User", entity.UserRoleStudent)
	repo.Create(context.Background(), user1)
	user2, _ := entity.NewUser("user2@example.com", "pass", "Second", "User", entity.UserRoleStudent)
	repo.Create(context.Background(), user2)

	users, err := repo.List(context.Background(), 10, 0)
	assert.NoError(t, err)
	assert.Len(t, users, 2)
}

func TestPostgresUserRepository_GetByRole(t *testing.T) {
	db := setupUserTestDB(t)
	repo := repository.NewPostgresUserRepository(db)

	student, _ := entity.NewUser("student@example.com", "pass", "Student", "User", entity.UserRoleStudent)
	repo.Create(context.Background(), student)
	instructor, _ := entity.NewUser("instructor@example.com", "pass", "Instructor", "User", entity.UserRoleInstructor)
	repo.Create(context.Background(), instructor)

	students, err := repo.GetByRole(context.Background(), entity.UserRoleStudent)
	assert.NoError(t, err)
	assert.Len(t, students, 1)
	assert.Equal(t, student.ID, students[0].ID)
}

func TestPostgresUserRepository_Get_NotFound(t *testing.T) {
	db := setupUserTestDB(t)
	repo := repository.NewPostgresUserRepository(db)

	t.Run("GetByID", func(t *testing.T) {
		_, err := repo.GetByID(context.Background(), uuid.New().String())
		assert.Error(t, err)
		assert.Equal(t, entity.ErrNotFound, err)
	})

	t.Run("GetByEmail", func(t *testing.T) {
		_, err := repo.GetByEmail(context.Background(), "nonexistent@example.com")
		assert.Error(t, err)
		assert.Equal(t, entity.ErrNotFound, err)
	})
}

func TestPostgresUserRepository_Create_DuplicateEmail(t *testing.T) {
	db := setupUserTestDB(t)
	repo := repository.NewPostgresUserRepository(db)

	user, err := entity.NewUser("test@example.com", "password123", "Test", "User", entity.UserRoleStudent)
	assert.NoError(t, err)
	err = repo.Create(context.Background(), user)
	assert.NoError(t, err)

	duplicateUser, err := entity.NewUser("test@example.com", "password456", "Another", "User", entity.UserRoleStudent)
	assert.NoError(t, err)
	err = repo.Create(context.Background(), duplicateUser)
	assert.Error(t, err) // Should be a unique constraint violation
}
