package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/config"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	domainRepository "github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/repository"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/infrastructure/database"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/infrastructure/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupEnrollmentTestDB(t *testing.T) (*gorm.DB, *entity.User, *entity.Course) {
	t.Helper()

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Skipf("skipping enrollment repo tests: %v", err)
	}

	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		t.Skipf("skipping enrollment repo tests: %v", err)
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

	student, err := entity.NewUser("student@example.com", "password123", "Student", "User", entity.UserRoleStudent)
	assert.NoError(t, err)
	err = gormDB.Create(student).Error
	assert.NoError(t, err)

	course, err := entity.NewCourse("Test Course", "Test Description", instructor.ID, entity.DifficultyLevelBeginner, instructor)
	assert.NoError(t, err)
	err = gormDB.Create(course).Error
	assert.NoError(t, err)

	return gormDB, student, course
}

func TestNewPostgresEnrollmentRepository(t *testing.T) {
	db, _, _ := setupEnrollmentTestDB(t)
	repo := repository.NewPostgresEnrollmentRepository(db)
	assert.NotNil(t, repo)
}

func TestPostgresEnrollmentRepository_CreateAndGet(t *testing.T) {
	db, student, course := setupEnrollmentTestDB(t)
	repo := repository.NewPostgresEnrollmentRepository(db)

	enrollment, err := entity.NewEnrollment(student, course.ID)
	assert.NoError(t, err)

	err = repo.Create(context.Background(), enrollment)
	assert.NoError(t, err)

	found, err := repo.GetByStudentAndCourse(context.Background(), student.ID, course.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, enrollment.StudentID, found.StudentID)
	assert.Equal(t, enrollment.CourseID, found.CourseID)
}

func TestPostgresEnrollmentRepository_Update(t *testing.T) {
	db, student, course := setupEnrollmentTestDB(t)
	repo := repository.NewPostgresEnrollmentRepository(db)

	enrollment, err := entity.NewEnrollment(student, course.ID)
	assert.NoError(t, err)
	err = repo.Create(context.Background(), enrollment)
	assert.NoError(t, err)

	enrollment.Status = entity.EnrollmentStatusCompleted
	err = repo.Update(context.Background(), enrollment)
	assert.NoError(t, err)

	found, err := repo.GetByStudentAndCourse(context.Background(), student.ID, course.ID)
	assert.NoError(t, err)
	assert.Equal(t, entity.EnrollmentStatusCompleted, found.Status)
}

func TestPostgresEnrollmentRepository_Delete(t *testing.T) {
	db, student, course := setupEnrollmentTestDB(t)
	repo := repository.NewPostgresEnrollmentRepository(db)

	enrollment, err := entity.NewEnrollment(student, course.ID)
	assert.NoError(t, err)
	err = repo.Create(context.Background(), enrollment)
	assert.NoError(t, err)

	err = repo.Delete(context.Background(), student.ID, course.ID)
	assert.NoError(t, err)

	_, err = repo.GetByStudentAndCourse(context.Background(), student.ID, course.ID)
	assert.Error(t, err)
	assert.Equal(t, entity.ErrNotFound, err)
}

func TestPostgresEnrollmentRepository_GetByStudent(t *testing.T) {
	db, student, course := setupEnrollmentTestDB(t)
	repo := repository.NewPostgresEnrollmentRepository(db)

	enrollment, err := entity.NewEnrollment(student, course.ID)
	assert.NoError(t, err)
	err = repo.Create(context.Background(), enrollment)
	assert.NoError(t, err)

	enrollments, total, err := repo.GetByStudent(context.Background(), student.ID, nil)
	assert.NoError(t, err)
	assert.Len(t, enrollments, 1)
	assert.Equal(t, int64(1), total)
}

func TestPostgresEnrollmentRepository_GetByCourse(t *testing.T) {
	db, student, course := setupEnrollmentTestDB(t)
	repo := repository.NewPostgresEnrollmentRepository(db)

	enrollment, err := entity.NewEnrollment(student, course.ID)
	assert.NoError(t, err)
	err = repo.Create(context.Background(), enrollment)
	assert.NoError(t, err)

	enrollments, total, err := repo.GetByCourse(context.Background(), course.ID, nil)
	assert.NoError(t, err)
	assert.Len(t, enrollments, 1)
	assert.Equal(t, int64(1), total)
}

func TestPostgresEnrollmentRepository_GetByStudent_WithFilters(t *testing.T) {
	db, student, course := setupEnrollmentTestDB(t)
	repo := repository.NewPostgresEnrollmentRepository(db)

	enrollment1, err := entity.NewEnrollment(student, course.ID)
	assert.NoError(t, err)
	err = repo.Create(context.Background(), enrollment1)
	assert.NoError(t, err)

	time.Sleep(1 * time.Second) // ensure different timestamps

	course2, err := entity.NewCourse("Another Course", "Desc", course.InstructorID, entity.DifficultyLevelIntermediate, nil)
	assert.NoError(t, err)
	db.Create(course2)

	enrollment2, err := entity.NewEnrollment(student, course2.ID)
	assert.NoError(t, err)
	enrollment2.Status = entity.EnrollmentStatusCompleted
	err = repo.Create(context.Background(), enrollment2)
	assert.NoError(t, err)

	t.Run("Filter by Status", func(t *testing.T) {
		status := entity.EnrollmentStatusActive
		enrollments, total, err := repo.GetByStudent(context.Background(), student.ID, &domainRepository.EnrollmentFilter{Status: &status})
		assert.NoError(t, err)
		assert.Len(t, enrollments, 1)
		assert.Equal(t, enrollment1.CourseID, enrollments[0].CourseID)
		assert.Equal(t, int64(1), total)
	})

	t.Run("Filter by DateFrom", func(t *testing.T) {
		dateFrom := time.Now().Add(-5 * time.Second).Format("2006-01-02")
		enrollments, total, err := repo.GetByStudent(context.Background(), student.ID, &domainRepository.EnrollmentFilter{DateFrom: &dateFrom})
		assert.NoError(t, err)
		assert.Len(t, enrollments, 2)
		assert.Equal(t, int64(2), total)
	})

	t.Run("Filter by DateTo", func(t *testing.T) {
		dateTo := time.Now().Add(-2 * time.Second).Format("2006-01-02")
		enrollments, total, err := repo.GetByStudent(context.Background(), student.ID, &domainRepository.EnrollmentFilter{DateTo: &dateTo})
		assert.NoError(t, err)
		assert.Len(t, enrollments, 0)
		assert.Equal(t, int64(0), total)
	})

	t.Run("Filter with Pagination", func(t *testing.T) {
		enrollments, total, err := repo.GetByStudent(context.Background(), student.ID, &domainRepository.EnrollmentFilter{Limit: 1, Offset: 0})
		assert.NoError(t, err)
		assert.Len(t, enrollments, 1)
		assert.Equal(t, enrollment2.CourseID, enrollments[0].CourseID) // Ordered by date DESC
		assert.Equal(t, int64(2), total)
	})
}

func TestPostgresEnrollmentRepository_GetByStudentAndCourse_NotFound(t *testing.T) {
	db, _, _ := setupEnrollmentTestDB(t)
	repo := repository.NewPostgresEnrollmentRepository(db)

	_, err := repo.GetByStudentAndCourse(context.Background(), uuid.New(), uuid.New())
	assert.Error(t, err)
	assert.Equal(t, entity.ErrNotFound, err)
}
