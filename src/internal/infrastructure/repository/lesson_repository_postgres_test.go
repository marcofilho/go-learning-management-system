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
	if err != nil {
		t.Skipf("skipping lesson repo tests: %v", err)
	}

	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		t.Skipf("skipping lesson repo tests: %v", err)
	}

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

func TestPostgresLessonRepository_CreateAndGetVersionByID(t *testing.T) {
	db, module := setupLessonTestDB(t)
	repo := repository.NewPostgresLessonRepository(db)

	lessonThread, err := entity.NewLesson(module.ID)
	assert.NoError(t, err)
	err = repo.CreateLesson(context.Background(), lessonThread)
	assert.NoError(t, err)

	lessonVersion := &entity.LessonVersion{
		LessonID:      lessonThread.ID,
		ModuleID:      module.ID,
		VersionNumber: 1,
		Content:       "Content",
	}
	err = repo.CreateVersion(context.Background(), lessonVersion)
	assert.NoError(t, err)

	found, err := repo.GetVersionByID(context.Background(), lessonVersion.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, lessonVersion.ID, found.ID)
}

func TestPostgresLessonRepository_GetLatestByModule(t *testing.T) {
	db, module := setupLessonTestDB(t)
	repo := repository.NewPostgresLessonRepository(db)

	lessonThread, err := entity.NewLesson(module.ID)
	assert.NoError(t, err)
	err = repo.CreateLesson(context.Background(), lessonThread)
	assert.NoError(t, err)

	lessonV1 := &entity.LessonVersion{LessonID: lessonThread.ID, ModuleID: module.ID, VersionNumber: 1, Content: "Content 1"}
	lessonV2 := &entity.LessonVersion{LessonID: lessonThread.ID, ModuleID: module.ID, VersionNumber: 2, Content: "Content 2"}

	assert.NoError(t, repo.CreateVersion(context.Background(), lessonV1))
	assert.NoError(t, repo.CreateVersion(context.Background(), lessonV2))

	lessons, total, err := repo.GetLatestByModule(context.Background(), module.ID, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, lessons, 1)
	assert.Equal(t, lessonV2.ID, lessons[0].ID)
	assert.Equal(t, int64(1), total)
}

func TestPostgresLessonRepository_GetAllVersions(t *testing.T) {
	db, module := setupLessonTestDB(t)
	repo := repository.NewPostgresLessonRepository(db)

	lessonThread, err := entity.NewLesson(module.ID)
	assert.NoError(t, err)
	assert.NoError(t, repo.CreateLesson(context.Background(), lessonThread))

	lessonV1 := &entity.LessonVersion{LessonID: lessonThread.ID, ModuleID: module.ID, VersionNumber: 1, Content: "Content 1"}
	lessonV2 := &entity.LessonVersion{LessonID: lessonThread.ID, ModuleID: module.ID, VersionNumber: 2, Content: "Content 2"}
	assert.NoError(t, repo.CreateVersion(context.Background(), lessonV1))
	assert.NoError(t, repo.CreateVersion(context.Background(), lessonV2))

	versions, total, err := repo.GetAllVersions(context.Background(), lessonThread.ID, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, versions, 2)
	assert.Equal(t, int64(2), total)
}

func TestPostgresLessonRepository_GetNextVersionNumber(t *testing.T) {
	db, module := setupLessonTestDB(t)
	repo := repository.NewPostgresLessonRepository(db)

	lessonThread, err := entity.NewLesson(module.ID)
	assert.NoError(t, err)
	assert.NoError(t, repo.CreateLesson(context.Background(), lessonThread))

	lessonV1 := &entity.LessonVersion{LessonID: lessonThread.ID, ModuleID: module.ID, VersionNumber: 1, Content: "Content"}
	assert.NoError(t, repo.CreateVersion(context.Background(), lessonV1))

	nextVersion, err := repo.GetNextVersionNumber(context.Background(), lessonThread.ID)
	assert.NoError(t, err)
	assert.Equal(t, 2, nextVersion)
}

func TestPostgresLessonRepository_Delete(t *testing.T) {
	db, module := setupLessonTestDB(t)
	repo := repository.NewPostgresLessonRepository(db)

	lessonThread, err := entity.NewLesson(module.ID)
	assert.NoError(t, err)
	assert.NoError(t, repo.CreateLesson(context.Background(), lessonThread))

	lessonV1 := &entity.LessonVersion{LessonID: lessonThread.ID, ModuleID: module.ID, VersionNumber: 1, Content: "Content"}
	assert.NoError(t, repo.CreateVersion(context.Background(), lessonV1))

	assert.NoError(t, repo.DeleteVersionsByLesson(context.Background(), lessonThread.ID))
	assert.NoError(t, repo.DeleteLesson(context.Background(), lessonThread.ID))

	_, err = repo.GetLessonByID(context.Background(), lessonThread.ID)
	assert.Error(t, err)
	assert.Equal(t, entity.ErrNotFound, err)
}

func TestPostgresLessonRepository_GetLessonByID_NotFound(t *testing.T) {
	db, _ := setupLessonTestDB(t)
	repo := repository.NewPostgresLessonRepository(db)

	_, err := repo.GetLessonByID(context.Background(), uuid.New())
	assert.Error(t, err)
	assert.Equal(t, entity.ErrNotFound, err)
}

func TestPostgresLessonRepository_DeleteLesson_NotFound(t *testing.T) {
	db, _ := setupLessonTestDB(t)
	repo := repository.NewPostgresLessonRepository(db)

	err := repo.DeleteLesson(context.Background(), uuid.New())
	assert.NoError(t, err)
}

func TestPostgresLessonRepository_GetLatestByModule_Empty(t *testing.T) {
	db, _ := setupLessonTestDB(t)
	repo := repository.NewPostgresLessonRepository(db)

	lessons, total, err := repo.GetLatestByModule(context.Background(), uuid.New(), 10, 0)
	assert.NoError(t, err)
	assert.Len(t, lessons, 0)
	assert.Equal(t, int64(0), total)
}

func TestPostgresLessonRepository_GetNextVersionNumber_Empty(t *testing.T) {
	db, _ := setupLessonTestDB(t)
	repo := repository.NewPostgresLessonRepository(db)

	nextVersion, err := repo.GetNextVersionNumber(context.Background(), uuid.New())
	assert.NoError(t, err)
	assert.Equal(t, 1, nextVersion)
}
