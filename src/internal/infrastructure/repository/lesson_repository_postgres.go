package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/repository"
)

type PostgresLessonRepository struct {
	db *gorm.DB
}

func NewPostgresLessonRepository(db *gorm.DB) repository.LessonRepository {
	return &PostgresLessonRepository{db: db}
}

func (r *PostgresLessonRepository) CreateLesson(ctx context.Context, lesson *entity.Lesson) error {
	return r.db.WithContext(ctx).Create(lesson).Error
}

func (r *PostgresLessonRepository) GetLessonByID(ctx context.Context, id uuid.UUID) (*entity.Lesson, error) {
	var lesson entity.Lesson
	if err := r.db.WithContext(ctx).First(&lesson, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}
	return &lesson, nil
}

func (r *PostgresLessonRepository) DeleteLesson(ctx context.Context, lessonID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Lesson{}, "id = ?", lessonID).Error
}

func (r *PostgresLessonRepository) CreateVersion(ctx context.Context, lesson *entity.LessonVersion) error {
	return r.db.WithContext(ctx).Create(lesson).Error
}

func (r *PostgresLessonRepository) GetVersionByID(ctx context.Context, id uuid.UUID) (*entity.LessonVersion, error) {
	var lesson entity.LessonVersion
	if err := r.db.WithContext(ctx).First(&lesson, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}
	return &lesson, nil
}

func (r *PostgresLessonRepository) GetLatestByModule(ctx context.Context, moduleID uuid.UUID, limit, offset int) ([]*entity.LessonVersion, int64, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, 0, err
	}

	var total int64
	countSQL := `SELECT COUNT(*) FROM lesson_versions lv
		INNER JOIN lessons l ON l.id = lv.lesson_id
		INNER JOIN (
			SELECT lv2.lesson_id, MAX(lv2.version_number) AS max_version
			FROM lesson_versions lv2
			INNER JOIN lessons l2 ON l2.id = lv2.lesson_id
			WHERE l2.module_id = $1
			GROUP BY lv2.lesson_id
		) latest ON latest.lesson_id = lv.lesson_id AND latest.max_version = lv.version_number
		WHERE l.module_id = $2`
	err = sqlDB.QueryRowContext(ctx, countSQL, moduleID, moduleID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	querySQL := `SELECT lv.id, lv.lesson_id, lv.version_number, lv.content, lv.video_url, lv.created_at, l.module_id FROM lesson_versions lv
		INNER JOIN lessons l ON l.id = lv.lesson_id
		INNER JOIN (
			SELECT lv2.lesson_id, MAX(lv2.version_number) AS max_version
			FROM lesson_versions lv2
			INNER JOIN lessons l2 ON l2.id = lv2.lesson_id
			WHERE l2.module_id = $1
			GROUP BY lv2.lesson_id
		) latest ON latest.lesson_id = lv.lesson_id AND latest.max_version = lv.version_number
		WHERE l.module_id = $2
		ORDER BY lv.created_at DESC`
	
	args := []interface{}{moduleID, moduleID}
	argIndex := 3
	if limit > 0 {
		querySQL += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, limit)
		argIndex++
		if offset > 0 {
			querySQL += fmt.Sprintf(" OFFSET $%d", argIndex)
			args = append(args, offset)
		}
	}

	rows, err := sqlDB.QueryContext(ctx, querySQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	type lessonResult struct {
		ID            uuid.UUID
		LessonID      uuid.UUID
		VersionNumber int
		Content       sql.NullString
		VideoURL      sql.NullString
		CreatedAt     time.Time
		ModuleID      uuid.UUID
	}
	var results []lessonResult
	for rows.Next() {
		var result lessonResult
		if err := rows.Scan(&result.ID, &result.LessonID, &result.VersionNumber, &result.Content, &result.VideoURL, &result.CreatedAt, &result.ModuleID); err != nil {
			return nil, 0, err
		}
		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	var lessons []*entity.LessonVersion
	lessons = make([]*entity.LessonVersion, len(results))
	for i, result := range results {
		lessons[i] = &entity.LessonVersion{
			ID:            result.ID,
			LessonID:      result.LessonID,
			ModuleID:      result.ModuleID,
			VersionNumber: result.VersionNumber,
			Content:       result.Content.String,
			VideoURL:      result.VideoURL.String,
			CreatedAt:     result.CreatedAt,
		}
	}

	return lessons, total, nil
}

func (r *PostgresLessonRepository) GetAllVersions(ctx context.Context, lessonID uuid.UUID, limit, offset int) ([]*entity.LessonVersion, int64, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, 0, err
	}

	var total int64
	err = sqlDB.QueryRowContext(ctx, "SELECT COUNT(*) FROM lesson_versions WHERE lesson_id = $1", lessonID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	querySQL := "SELECT lv.id, lv.lesson_id, lv.version_number, lv.content, lv.video_url, lv.created_at, l.module_id FROM lesson_versions lv INNER JOIN lessons l ON l.id = lv.lesson_id WHERE lv.lesson_id = $1 ORDER BY lv.version_number DESC"
	
	args := []interface{}{lessonID}
	argIndex := 2
	if limit > 0 {
		querySQL += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, limit)
		argIndex++
		if offset > 0 {
			querySQL += fmt.Sprintf(" OFFSET $%d", argIndex)
			args = append(args, offset)
		}
	}

	rows, err := sqlDB.QueryContext(ctx, querySQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	type versionResult struct {
		ID            uuid.UUID
		LessonID      uuid.UUID
		VersionNumber int
		Content       sql.NullString
		VideoURL      sql.NullString
		CreatedAt     time.Time
		ModuleID      uuid.UUID
	}
	var results []versionResult
	for rows.Next() {
		var result versionResult
		if err := rows.Scan(&result.ID, &result.LessonID, &result.VersionNumber, &result.Content, &result.VideoURL, &result.CreatedAt, &result.ModuleID); err != nil {
			return nil, 0, err
		}
		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	var versions []*entity.LessonVersion

	versions = make([]*entity.LessonVersion, len(results))
	for i, result := range results {
		versions[i] = &entity.LessonVersion{
			ID:            result.ID,
			LessonID:      result.LessonID,
			ModuleID:      result.ModuleID,
			VersionNumber: result.VersionNumber,
			Content:       result.Content.String,
			VideoURL:      result.VideoURL.String,
			CreatedAt:     result.CreatedAt,
		}
	}

	return versions, total, nil
}

func (r *PostgresLessonRepository) GetNextVersionNumber(ctx context.Context, lessonID uuid.UUID) (int, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return 0, err
	}

	var maxVersion int
	err = sqlDB.QueryRowContext(ctx, "SELECT COALESCE(MAX(version_number), 0) FROM lesson_versions WHERE lesson_id = $1", lessonID).Scan(&maxVersion)
	if err != nil {
		return 0, err
	}
	return maxVersion + 1, nil
}

func (r *PostgresLessonRepository) DeleteVersionsByLesson(ctx context.Context, lessonID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("lesson_id = ?", lessonID).Delete(&entity.LessonVersion{}).Error
}
