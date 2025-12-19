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

func setupModuleTestDB(t *testing.T) (*gorm.DB, *entity.Course) {
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

	instructor, err := entity.NewUser("instructor@example.com", "password123", "Instructor", "User", entity.UserRoleInstructor)
	assert.NoError(t, err)
	err = gormDB.Create(instructor).Error
	assert.NoError(t, err)

	course, err := entity.NewCourse("Test Course", "Test Description", instructor.ID, entity.DifficultyLevelBeginner, instructor)
	assert.NoError(t, err)
	err = gormDB.Create(course).Error
	assert.NoError(t, err)

	return gormDB, course
}

func TestNewPostgresModuleRepository(t *testing.T) {
	db, _ := setupModuleTestDB(t)
	repo := repository.NewPostgresModuleRepository(db)
	assert.NotNil(t, repo)
}

func TestPostgresModuleRepository_CreateAndGetByID(t *testing.T) {
	db, course := setupModuleTestDB(t)
	repo := repository.NewPostgresModuleRepository(db)

	module, err := entity.NewModule("Test Module", course.ID)
	assert.NoError(t, err)

	err = repo.Create(context.Background(), module)
	assert.NoError(t, err)

	found, err := repo.GetByID(context.Background(), module.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, module.ID, found.ID)
}

func TestPostgresModuleRepository_Update(t *testing.T) {
	db, course := setupModuleTestDB(t)
	repo := repository.NewPostgresModuleRepository(db)

	module, err := entity.NewModule("Test Module", course.ID)
	assert.NoError(t, err)
	err = repo.Create(context.Background(), module)
	assert.NoError(t, err)

	module.Title = "Updated Title"
	err = repo.Update(context.Background(), module)
	assert.NoError(t, err)

	found, err := repo.GetByID(context.Background(), module.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", found.Title)
}

func TestPostgresModuleRepository_Delete(t *testing.T) {
	db, course := setupModuleTestDB(t)
	repo := repository.NewPostgresModuleRepository(db)

	module, err := entity.NewModule("Test Module", course.ID)
	assert.NoError(t, err)
	err = repo.Create(context.Background(), module)
	assert.NoError(t, err)

	err = repo.Delete(context.Background(), module.ID)
	assert.NoError(t, err)

	_, err = repo.GetByID(context.Background(), module.ID)
	assert.Error(t, err)
	assert.Equal(t, entity.ErrNotFound, err)
}

func TestPostgresModuleRepository_GetByCourse(t *testing.T) {
	db, course := setupModuleTestDB(t)
	repo := repository.NewPostgresModuleRepository(db)

	module1, err := entity.NewModule("Module 1", course.ID)
	assert.NoError(t, err)
	err = repo.Create(context.Background(), module1)
	assert.NoError(t, err)

	module2, err := entity.NewModule("Module 2", course.ID)
	assert.NoError(t, err)
	err = repo.Create(context.Background(), module2)
	assert.NoError(t, err)

	modules, err := repo.GetByCourse(context.Background(), course.ID)
	assert.NoError(t, err)
	assert.Len(t, modules, 2)
}

func TestPostgresModuleRepository_GetByID_NotFound(t *testing.T) {
	db, _ := setupModuleTestDB(t)
	repo := repository.NewPostgresModuleRepository(db)

	_, err := repo.GetByID(context.Background(), uuid.New().String())
	assert.Error(t, err)
	assert.Equal(t, entity.ErrNotFound, err)
}

func TestPostgresModuleRepository_Delete_NotFound(t *testing.T) {
	db, _ := setupModuleTestDB(t)
	repo := repository.NewPostgresModuleRepository(db)

	err := repo.Delete(context.Background(), uuid.New().String())
	assert.NoError(t, err)
}

func TestPostgresModuleRepository_GetByCourse_Empty(t *testing.T) {
	db, _ := setupModuleTestDB(t)
	repo := repository.NewPostgresModuleRepository(db)

	modules, err := repo.GetByCourse(context.Background(), uuid.New().String())
	assert.NoError(t, err)
	assert.Len(t, modules, 0)
}
