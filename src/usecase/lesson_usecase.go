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

	lesson.VersionNumber = 1

	if err := lesson.Validate(); err != nil {
		return err
	}

	if err := uc.lessonRepo.Create(ctx, lesson); err != nil {
		return err
	}

	// Create audit log
	payloadAfter := fmt.Sprintf(`{"id":"%s","module_id":"%s","version":%d}`, lesson.ID, lesson.ModuleID, lesson.VersionNumber)
	auditLog, _ := entity.NewAuditLog(entity.AuditActionLessonVersionCreated, lesson.ID, "lesson", "", payloadAfter, &instructorID)
	uc.auditLogRepo.Create(ctx, auditLog)

	return nil
}

func (uc *LessonUseCase) CreateLessonVersion(ctx context.Context, lessonID uuid.UUID, newVersion *entity.LessonVersion, instructorID uuid.UUID) error {
	existingLesson, err := uc.lessonRepo.GetByID(ctx, lessonID)
	if err != nil {
		return err
	}

	module, err := uc.moduleRepo.GetByID(ctx, existingLesson.ModuleID)
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

	nextVersion, err := uc.lessonRepo.GetNextVersionNumber(ctx, existingLesson.ModuleID)
	if err != nil {
		return err
	}

	newVersion.ModuleID = existingLesson.ModuleID
	newVersion.VersionNumber = nextVersion

	if err := newVersion.Validate(); err != nil {
		return err
	}

	if err := uc.lessonRepo.Create(ctx, newVersion); err != nil {
		return err
	}

	// Create audit log
	payloadBefore := fmt.Sprintf(`{"id":"%s","version":%d}`, existingLesson.ID, existingLesson.VersionNumber)
	payloadAfter := fmt.Sprintf(`{"id":"%s","version":%d}`, newVersion.ID, newVersion.VersionNumber)
	auditLog, _ := entity.NewAuditLog(entity.AuditActionLessonVersionCreated, newVersion.ID, "lesson", payloadBefore, payloadAfter, &instructorID)
	uc.auditLogRepo.Create(ctx, auditLog)

	return nil
}

func (uc *LessonUseCase) GetLatestLessonsByModule(ctx context.Context, moduleID uuid.UUID) ([]*entity.LessonVersion, error) {
	if _, err := uc.moduleRepo.GetByID(ctx, moduleID); err != nil {
		return nil, err
	}

	return uc.lessonRepo.GetLatestByModule(ctx, moduleID)
}

func (uc *LessonUseCase) GetAllLessonVersions(ctx context.Context, lessonID uuid.UUID) ([]*entity.LessonVersion, error) {
	return uc.lessonRepo.GetAllVersions(ctx, lessonID)
}

func (uc *LessonUseCase) GetLesson(ctx context.Context, id uuid.UUID) (*entity.LessonVersion, error) {
	return uc.lessonRepo.GetByID(ctx, id)
}

func (uc *LessonUseCase) DeleteLesson(ctx context.Context, id uuid.UUID, instructorID uuid.UUID) error {
	lesson, err := uc.lessonRepo.GetByID(ctx, id)
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

	if err := uc.lessonRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Create audit log
	payloadBefore := fmt.Sprintf(`{"id":"%s","module_id":"%s"}`, lesson.ID, lesson.ModuleID)
	auditLog, _ := entity.NewAuditLog(entity.AuditActionLessonDeleted, id, "lesson", payloadBefore, "", &instructorID)
	uc.auditLogRepo.Create(ctx, auditLog)

	return nil
}
