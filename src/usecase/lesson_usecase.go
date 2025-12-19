package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/repository"
)

type LessonUseCase struct {
	lessonRepo   repository.LessonRepository
	moduleRepo   repository.ModuleRepository
	courseRepo   repository.CourseRepository
	auditLogRepo repository.AuditLogRepository
}

func NewLessonUseCase(lessonRepo repository.LessonRepository, moduleRepo repository.ModuleRepository, courseRepo repository.CourseRepository, auditLogRepo repository.AuditLogRepository) *LessonUseCase {
	return &LessonUseCase{
		lessonRepo:   lessonRepo,
		moduleRepo:   moduleRepo,
		courseRepo:   courseRepo,
		auditLogRepo: auditLogRepo,
	}
}

func (uc *LessonUseCase) CreateLesson(ctx context.Context, lesson *entity.LessonVersion, instructorID uuid.UUID) error {
	module, err := uc.moduleRepo.GetByID(ctx, lesson.ModuleID)
	if err != nil {
		return err
	}

	course, err := uc.courseRepo.GetByID(ctx, module.CourseID)
	if err != nil {
		return err
	}

	if !lesson.CanBeModifiedBy(course, instructorID) {
		return entity.ErrUnauthorized
	}

	// Create a lesson thread first
	baseLesson, err := entity.NewLesson(module.ID)
	if err != nil {
		return err
	}
	if err := uc.lessonRepo.CreateLesson(ctx, baseLesson); err != nil {
		return err
	}

	lesson.LessonID = baseLesson.ID
	lesson.ModuleID = module.ID
	lesson.VersionNumber = 1

	if err := lesson.Validate(); err != nil {
		return err
	}

	if err := uc.lessonRepo.CreateVersion(ctx, lesson); err != nil {
		return err
	}

	// Create audit log
	payloadAfter := fmt.Sprintf(`{"id":"%s","module_id":"%s","version":%d}`, lesson.ID, lesson.ModuleID, lesson.VersionNumber)
	auditLog, _ := entity.NewAuditLog(entity.AuditActionLessonVersionCreated, lesson.ID, "lesson", "", payloadAfter, &instructorID)
	uc.auditLogRepo.Create(ctx, auditLog)

	return nil
}

func (uc *LessonUseCase) CreateLessonVersion(ctx context.Context, lessonID uuid.UUID, newVersion *entity.LessonVersion, instructorID uuid.UUID) error {
	baseLesson, err := uc.lessonRepo.GetLessonByID(ctx, lessonID)
	if err != nil {
		return err
	}

	module, err := uc.moduleRepo.GetByID(ctx, baseLesson.ModuleID)
	if err != nil {
		return err
	}

	course, err := uc.courseRepo.GetByID(ctx, module.CourseID)
	if err != nil {
		return err
	}

	if !newVersion.CanBeModifiedBy(course, instructorID) {
		return entity.ErrUnauthorized
	}

	nextVersion, err := uc.lessonRepo.GetNextVersionNumber(ctx, lessonID)
	if err != nil {
		return err
	}

	newVersion.LessonID = lessonID
	newVersion.ModuleID = baseLesson.ModuleID
	newVersion.VersionNumber = nextVersion

	if err := newVersion.Validate(); err != nil {
		return err
	}

	if err := uc.lessonRepo.CreateVersion(ctx, newVersion); err != nil {
		return err
	}

	// Create audit log
	payloadBefore := fmt.Sprintf(`{"lesson_id":"%s","version":%d}`, lessonID, nextVersion-1)
	payloadAfter := fmt.Sprintf(`{"lesson_id":"%s","version":%d}`, lessonID, newVersion.VersionNumber)
	auditLog, _ := entity.NewAuditLog(entity.AuditActionLessonVersionCreated, newVersion.ID, "lesson", payloadBefore, payloadAfter, &instructorID)
	uc.auditLogRepo.Create(ctx, auditLog)

	return nil
}

func (uc *LessonUseCase) GetLatestLessonsByModule(ctx context.Context, moduleID uuid.UUID, limit, offset int) ([]*entity.LessonVersion, int64, error) {
	if _, err := uc.moduleRepo.GetByID(ctx, moduleID); err != nil {
		return nil, 0, err
	}

	return uc.lessonRepo.GetLatestByModule(ctx, moduleID, limit, offset)
}

func (uc *LessonUseCase) GetAllLessonVersions(ctx context.Context, lessonID uuid.UUID, limit, offset int) ([]*entity.LessonVersion, int64, error) {
	return uc.lessonRepo.GetAllVersions(ctx, lessonID, limit, offset)
}

func (uc *LessonUseCase) GetLesson(ctx context.Context, id uuid.UUID) (*entity.LessonVersion, error) {
	return uc.lessonRepo.GetVersionByID(ctx, id)
}

func (uc *LessonUseCase) DeleteLesson(ctx context.Context, id uuid.UUID, instructorID uuid.UUID) error {
	lesson, err := uc.lessonRepo.GetLessonByID(ctx, id)
	if err != nil {
		return err
	}

	module, err := uc.moduleRepo.GetByID(ctx, lesson.ModuleID)
	if err != nil {
		return err
	}

	course, err := uc.courseRepo.GetByID(ctx, module.CourseID)
	if err != nil {
		return err
	}

	if !lesson.CanBeModifiedBy(course, instructorID) {
		return entity.ErrUnauthorized
	}

	if err := uc.lessonRepo.DeleteVersionsByLesson(ctx, id); err != nil {
		return err
	}

	if err := uc.lessonRepo.DeleteLesson(ctx, id); err != nil {
		return err
	}

	// Create audit log
	payloadBefore := fmt.Sprintf(`{"lesson_id":"%s","module_id":"%s"}`, lesson.ID, lesson.ModuleID)
	auditLog, _ := entity.NewAuditLog(entity.AuditActionLessonDeleted, id, "lesson", payloadBefore, "", &instructorID)
	uc.auditLogRepo.Create(ctx, auditLog)

	return nil
}
