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
	// Verify module exists and instructor owns the course
	module, err := uc.moduleRepo.GetByID(ctx, lesson.ModuleID)
	if err != nil {
		return err
	}

	course, err := uc.courseRepo.GetByID(ctx, module.CourseID)
	if err != nil {
		return err
	}

	if course.InstructorID != instructorID {
		return entity.ErrUnauthorized
	}

	// Auto-set version number to 1 for new lessons
	lesson.VersionNumber = 1

	return uc.lessonRepo.Create(ctx, lesson)
}

func (uc *LessonUseCase) CreateLessonVersion(ctx context.Context, lessonID string, newVersion *entity.LessonVersion, instructorID string) error {
	// Get existing lesson
	existingLesson, err := uc.lessonRepo.GetByID(ctx, lessonID)
	if err != nil {
		return err
	}

	// Verify instructor owns the course
	module, err := uc.moduleRepo.GetByID(ctx, existingLesson.ModuleID)
	if err != nil {
		return err
	}

	course, err := uc.courseRepo.GetByID(ctx, module.CourseID)
	if err != nil {
		return err
	}

	if course.InstructorID != instructorID {
		return entity.ErrUnauthorized
	}

	// Get next version number
	nextVersion, err := uc.lessonRepo.GetNextVersionNumber(ctx, existingLesson.ModuleID)
	if err != nil {
		return err
	}

	// Set module ID and version number
	newVersion.ModuleID = existingLesson.ModuleID
	newVersion.VersionNumber = nextVersion

	return uc.lessonRepo.Create(ctx, newVersion)
}

func (uc *LessonUseCase) GetLatestLessonsByModule(ctx context.Context, moduleID string) ([]*entity.LessonVersion, error) {
	// Verify module exists
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
	// Get lesson to check ownership
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

	if course.InstructorID != instructorID {
		return entity.ErrUnauthorized
	}

	return uc.lessonRepo.Delete(ctx, id)
}
