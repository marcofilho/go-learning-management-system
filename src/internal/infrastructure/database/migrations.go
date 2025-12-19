package database

import (
	"fmt"

	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return fmt.Errorf("failed to create uuid extension: %w", err)
	}

	if err := db.AutoMigrate(
		&entity.User{},
		&entity.Course{},
		&entity.Module{},
		&entity.LessonVersion{},
		&entity.Enrollment{},
		&entity.AuditLog{},
	); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if err := createIndexes(db); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	return nil
}

func createIndexes(db *gorm.DB) error {
	// User indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_users_role ON users(role) WHERE deleted_at IS NULL").Error; err != nil {
		return err
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active) WHERE deleted_at IS NULL").Error; err != nil {
		return err
	}

	// Course indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_courses_instructor ON courses(instructor_id) WHERE deleted_at IS NULL").Error; err != nil {
		return err
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_courses_difficulty ON courses(difficulty_level) WHERE deleted_at IS NULL").Error; err != nil {
		return err
	}

	// Module indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_modules_course ON modules(course_id, order_index) WHERE deleted_at IS NULL").Error; err != nil {
		return err
	}

	// Lesson indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_lessons_module ON lesson_versions(module_id, version_number) WHERE deleted_at IS NULL").Error; err != nil {
		return err
	}

	// Enrollment indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_enrollments_student ON course_enrollments(student_id, status) WHERE deleted_at IS NULL").Error; err != nil {
		return err
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_enrollments_course ON course_enrollments(course_id, status) WHERE deleted_at IS NULL").Error; err != nil {
		return err
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_enrollments_date ON course_enrollments(enrollment_date) WHERE deleted_at IS NULL").Error; err != nil {
		return err
	}

	return nil
}

func DropAllTables(db *gorm.DB) error {
	return db.Migrator().DropTable(
		&entity.AuditLog{},
		&entity.Enrollment{},
		&entity.LessonVersion{},
		&entity.Module{},
		&entity.Course{},
		&entity.User{},
	)
}
