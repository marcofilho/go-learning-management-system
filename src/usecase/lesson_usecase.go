package usecase

import (
	"context"

	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/repository"
)

type LessonUseCase struct {
	lessonRepo repository.LessonRepository
	moduleRepo repository.ModuleRepository
	courseRepo repository.CourseRepository
}

func NewLessonUseCase(lessonRepo repository.LessonRepository, moduleRepo repository.ModuleRepository, courseRepo repository.CourseRepository) *LessonUseCase {
	return &LessonUseCase{
		lessonRepo: lessonRepo,
		moduleRepo: moduleRepo,
		courseRepo: courseRepo,
	}
}

func (uc *LessonUseCase) CreateLesson(ctx context.Context, lesson *entity.LessonVersion, instructorID string) error {
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

	return uc.lessonRepo.Create(ctx, lesson)
}

func (uc *LessonUseCase) CreateLessonVersion(ctx context.Context, lessonID string, newVersion *entity.LessonVersion, instructorID string) error {
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

	return uc.lessonRepo.Create(ctx, newVersion)
}

func (uc *LessonUseCase) GetLatestLessonsByModule(ctx context.Context, moduleID string) ([]*entity.LessonVersion, error) {
	if _, err := uc.moduleRepo.GetByID(ctx, moduleID); err != nil {
		return nil, err
	}

	return uc.lessonRepo.GetLatestByModule(ctx, moduleID)
}

func (uc *LessonUseCase) GetAllLessonVersions(ctx context.Context, lessonID string) ([]*entity.LessonVersion, error) {
	return uc.lessonRepo.GetAllVersions(ctx, lessonID)
}

func (uc *LessonUseCase) GetLesson(ctx context.Context, id string) (*entity.LessonVersion, error) {
	return uc.lessonRepo.GetByID(ctx, id)
}

func (uc *LessonUseCase) DeleteLesson(ctx context.Context, id string, instructorID string) error {
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

	return uc.lessonRepo.Delete(ctx, id)
}
