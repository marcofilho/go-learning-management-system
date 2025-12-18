package repository_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/config"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	domainRepository "github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/repository"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/infrastructure/database"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/infrastructure/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupCourseTestDB(t *testing.T) (*gorm.DB, *entity.User) {
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

	user, err := entity.NewUser("test@example.com", "password123", "Test", "User", entity.UserRoleInstructor)
	assert.NoError(t, err)
	err = gormDB.Create(user).Error
	assert.NoError(t, err)

	return gormDB, user
}

func TestNewPostgresCourseRepository(t *testing.T) {
	db, _ := setupCourseTestDB(t)
	repo := repository.NewPostgresCourseRepository(db)
	assert.NotNil(t, repo)
}

func TestPostgresCourseRepository_CreateAndGetByID(t *testing.T) {
	db, instructor := setupCourseTestDB(t)
	repo := repository.NewPostgresCourseRepository(db)

	course, err := entity.NewCourse("Test Course", "Test Description", instructor.ID, entity.DifficultyLevelBeginner, instructor)
	assert.NoError(t, err)

	err = repo.Create(context.Background(), course)
	assert.NoError(t, err)

	found, err := repo.GetByID(context.Background(), course.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, course.ID, found.ID)
}

func TestPostgresCourseRepository_Update(t *testing.T) {
	db, instructor := setupCourseTestDB(t)
	repo := repository.NewPostgresCourseRepository(db)

	course, err := entity.NewCourse("Test Course", "Test Description", instructor.ID, entity.DifficultyLevelBeginner, instructor)
	assert.NoError(t, err)

	err = repo.Create(context.Background(), course)
	assert.NoError(t, err)

	course.Title = "Updated Title"
	err = repo.Update(context.Background(), course)
	assert.NoError(t, err)

	found, err := repo.GetByID(context.Background(), course.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", found.Title)
}

func TestPostgresCourseRepository_Delete(t *testing.T) {
	db, instructor := setupCourseTestDB(t)
	repo := repository.NewPostgresCourseRepository(db)

	course, err := entity.NewCourse("Test Course", "Test Description", instructor.ID, entity.DifficultyLevelBeginner, instructor)
	assert.NoError(t, err)

	err = repo.Create(context.Background(), course)
	assert.NoError(t, err)

	err = repo.Delete(context.Background(), course.ID)
	assert.NoError(t, err)

	_, err = repo.GetByID(context.Background(), course.ID)
	assert.Error(t, err)
	assert.Equal(t, entity.ErrNotFound, err)
}

func TestPostgresCourseRepository_List(t *testing.T) {
	db, instructor := setupCourseTestDB(t)
	repo := repository.NewPostgresCourseRepository(db)

	course1, err := entity.NewCourse("Test Course 1", "Test Description", instructor.ID, entity.DifficultyLevelBeginner, instructor)
	assert.NoError(t, err)
	err = repo.Create(context.Background(), course1)
	assert.NoError(t, err)

	course2, err := entity.NewCourse("Test Course 2", "Test Description", instructor.ID, entity.DifficultyLevelBeginner, instructor)
	assert.NoError(t, err)
	err = repo.Create(context.Background(), course2)
	assert.NoError(t, err)

	courses, err := repo.List(context.Background(), &domainRepository.CourseFilter{Limit: 10, Offset: 0})
	assert.NoError(t, err)
	assert.NotNil(t, courses)
	assert.GreaterOrEqual(t, len(courses), 2)
}

func TestPostgresCourseRepository_GetByInstructor(t *testing.T) {
	db, instructor := setupCourseTestDB(t)
	repo := repository.NewPostgresCourseRepository(db)

	course1, err := entity.NewCourse("Test Course 1", "Test Description", instructor.ID, entity.DifficultyLevelBeginner, instructor)
	assert.NoError(t, err)
	err = repo.Create(context.Background(), course1)
	assert.NoError(t, err)

	course2, err := entity.NewCourse("Test Course 2", "Test Description", instructor.ID, entity.DifficultyLevelBeginner, instructor)
	assert.NoError(t, err)
	err = repo.Create(context.Background(), course2)
	assert.NoError(t, err)

	courses, err := repo.GetByInstructor(context.Background(), instructor.ID)
	assert.NoError(t, err)
	assert.NotNil(t, courses)
	assert.Len(t, courses, 2)
}

func TestPostgresCourseRepository_GetByInstructor_NotFound(t *testing.T) {
	db, _ := setupCourseTestDB(t)
	repo := repository.NewPostgresCourseRepository(db)

	courses, err := repo.GetByInstructor(context.Background(), uuid.New().String())
	assert.NoError(t, err)
	assert.NotNil(t, courses)
	assert.Len(t, courses, 0)
}

func TestPostgresCourseRepository_List_WithFilters(t *testing.T) {
	db, instructor := setupCourseTestDB(t)
	repo := repository.NewPostgresCourseRepository(db)

	course1, err := entity.NewCourse("Test Course 1", "Test Description", instructor.ID, entity.DifficultyLevelBeginner, instructor)
	assert.NoError(t, err)
	err = repo.Create(context.Background(), course1)
	assert.NoError(t, err)

	course2, err := entity.NewCourse("Test Course 2", "Test Description", instructor.ID, entity.DifficultyLevelIntermediate, instructor)
	assert.NoError(t, err)
	err = repo.Create(context.Background(), course2)
	assert.NoError(t, err)

	t.Run("Filter by InstructorID", func(t *testing.T) {
		courses, err := repo.List(context.Background(), &domainRepository.CourseFilter{InstructorID: &instructor.ID})
		assert.NoError(t, err)
		assert.Len(t, courses, 2)
	})

	t.Run("Filter by DifficultyLevel", func(t *testing.T) {
		difficulty := entity.DifficultyLevelBeginner
		courses, err := repo.List(context.Background(), &domainRepository.CourseFilter{DifficultyLevel: &difficulty})
		assert.NoError(t, err)
		assert.Len(t, courses, 1)
		assert.Equal(t, course1.ID, courses[0].ID)
	})

	t.Run("Filter by ActiveOnly", func(t *testing.T) {
		courses, err := repo.List(context.Background(), &domainRepository.CourseFilter{ActiveOnly: true})
		assert.NoError(t, err)
		assert.Len(t, courses, 2)
	})
}

func TestPostgresCourseRepository_GetByID_NotFound(t *testing.T) {
	db, _ := setupCourseTestDB(t)
	repo := repository.NewPostgresCourseRepository(db)

	_, err := repo.GetByID(context.Background(), uuid.New().String())
	assert.Error(t, err)
	assert.Equal(t, entity.ErrNotFound, err)
}
