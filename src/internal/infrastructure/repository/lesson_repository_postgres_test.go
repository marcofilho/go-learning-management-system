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

func setupLessonTestDB(t *testing.T) (*gorm.DB, *entity.Module) {
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

	module, err := entity.NewModule("Test Module", course.ID)
	assert.NoError(t, err)
	err = gormDB.Create(module).Error
	assert.NoError(t, err)

	return gormDB, module
}

func TestNewPostgresLessonRepository(t *testing.T) {
	db, _ := setupLessonTestDB(t)
	repo := repository.NewPostgresLessonRepository(db)
	assert.NotNil(t, repo)
}

func TestPostgresLessonRepository_CreateAndGetByID(t *testing.T) {
	db, module := setupLessonTestDB(t)
	repo := repository.NewPostgresLessonRepository(db)

	lesson, err := entity.NewLessonVersion("Test Lesson", "Content", module.ID, 1)
	assert.NoError(t, err)

	err = repo.Create(context.Background(), lesson)
	assert.NoError(t, err)

	found, err := repo.GetByID(context.Background(), lesson.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, lesson.ID, found.ID)
}

func TestPostgresLessonRepository_GetByModule(t *testing.T) {
	db, module := setupLessonTestDB(t)
	repo := repository.NewPostgresLessonRepository(db)

	lesson1, err := entity.NewLessonVersion("Lesson 1", "Content 1", module.ID, 1)
	assert.NoError(t, err)
	err = repo.Create(context.Background(), lesson1)
	assert.NoError(t, err)

	lesson2, err := entity.NewLessonVersion("Lesson 1", "Content 2", module.ID, 2)
	assert.NoError(t, err)
	err = repo.Create(context.Background(), lesson2)
	assert.NoError(t, err)

	lessons, err := repo.GetByModule(context.Background(), module.ID)
	assert.NoError(t, err)
	assert.Len(t, lessons, 2)
}

func TestPostgresLessonRepository_GetLatestByModule(t *testing.T) {
	db, module := setupLessonTestDB(t)
	repo := repository.NewPostgresLessonRepository(db)

	lesson1v1, err := entity.NewLessonVersion("Lesson 1", "Content 1", module.ID, 1)
	assert.NoError(t, err)
	err = repo.Create(context.Background(), lesson1v1)
	assert.NoError(t, err)

	lesson1v2, err := entity.NewLessonVersion("Lesson 1", "Content 2", module.ID, 2)
	assert.NoError(t, err)
	err = repo.Create(context.Background(), lesson1v2)
	assert.NoError(t, err)

	lessons, err := repo.GetLatestByModule(context.Background(), module.ID)
	assert.NoError(t, err)
	assert.Len(t, lessons, 1)
	assert.Equal(t, lesson1v2.ID, lessons[0].ID)
}

func TestPostgresLessonRepository_GetAllVersions(t *testing.T) {
	db, module := setupLessonTestDB(t)
	repo := repository.NewPostgresLessonRepository(db)

	lesson1v1, err := entity.NewLessonVersion("Lesson 1", "Content 1", module.ID, 1)
	assert.NoError(t, err)
	err = repo.Create(context.Background(), lesson1v1)
	assert.NoError(t, err)

	lesson1v2, err := entity.NewLessonVersion("Lesson 1", "Content 2", module.ID, 2)
	assert.NoError(t, err)
	err = repo.Create(context.Background(), lesson1v2)
	assert.NoError(t, err)

	versions, err := repo.GetAllVersions(context.Background(), lesson1v1.ID)
	assert.NoError(t, err)
	assert.Len(t, versions, 2)
}

func TestPostgresLessonRepository_GetNextVersionNumber(t *testing.T) {
	db, module := setupLessonTestDB(t)
	repo := repository.NewPostgresLessonRepository(db)

	lesson, err := entity.NewLessonVersion("Test Lesson", "Content", module.ID, 1)
	assert.NoError(t, err)
	err = repo.Create(context.Background(), lesson)
	assert.NoError(t, err)

	nextVersion, err := repo.GetNextVersionNumber(context.Background(), module.ID)
	assert.NoError(t, err)
	assert.Equal(t, 2, nextVersion)
}

func TestPostgresLessonRepository_Delete(t *testing.T) {
	db, module := setupLessonTestDB(t)
	repo := repository.NewPostgresLessonRepository(db)

	lesson, err := entity.NewLessonVersion("Test Lesson", "Content", module.ID, 1)
	assert.NoError(t, err)
	err = repo.Create(context.Background(), lesson)
	assert.NoError(t, err)

	err = repo.Delete(context.Background(), lesson.ID)
	assert.NoError(t, err)

	_, err = repo.GetByID(context.Background(), lesson.ID)
	assert.Error(t, err)
	assert.Equal(t, entity.ErrNotFound, err)
}

func TestPostgresLessonRepository_GetByID_NotFound(t *testing.T) {
	db, _ := setupLessonTestDB(t)
	repo := repository.NewPostgresLessonRepository(db)

	_, err := repo.GetByID(context.Background(), uuid.New())
	assert.Error(t, err)
	assert.Equal(t, entity.ErrNotFound, err)
}

func TestPostgresLessonRepository_Delete_NotFound(t *testing.T) {
	db, _ := setupLessonTestDB(t)
	repo := repository.NewPostgresLessonRepository(db)

	err := repo.Delete(context.Background(), uuid.New())
	assert.NoError(t, err)
}

func TestPostgresLessonRepository_GetByModule_Empty(t *testing.T) {
	db, _ := setupLessonTestDB(t)
	repo := repository.NewPostgresLessonRepository(db)

	lessons, err := repo.GetByModule(context.Background(), uuid.New())
	assert.NoError(t, err)
	assert.Len(t, lessons, 0)
}

func TestPostgresLessonRepository_GetNextVersionNumber_Empty(t *testing.T) {
	db, _ := setupLessonTestDB(t)
	repo := repository.NewPostgresLessonRepository(db)

	nextVersion, err := repo.GetNextVersionNumber(context.Background(), uuid.New())
	assert.NoError(t, err)
	assert.Equal(t, 1, nextVersion)
}
