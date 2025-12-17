package usecase

import (
	"context"

	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/repository"
)

type ModuleUseCase struct {
	moduleRepo repository.ModuleRepository
	courseRepo repository.CourseRepository
}

func NewModuleUseCase(moduleRepo repository.ModuleRepository, courseRepo repository.CourseRepository) *ModuleUseCase {
	return &ModuleUseCase{
		moduleRepo: moduleRepo,
		courseRepo: courseRepo,
	}
}

func (uc *ModuleUseCase) CreateModule(ctx context.Context, module *entity.Module, instructorID string) error {
	// Verify course exists and instructor owns it
	course, err := uc.courseRepo.GetByID(ctx, module.CourseID)
	if err != nil {
		return err
	}

	if course.InstructorID != instructorID {
		return entity.ErrUnauthorized
	}

	return uc.moduleRepo.Create(ctx, module)
}

func (uc *ModuleUseCase) GetModulesByCourse(ctx context.Context, courseID string) ([]*entity.Module, error) {
	// Verify course exists
	if _, err := uc.courseRepo.GetByID(ctx, courseID); err != nil {
		return nil, err
	}

	return uc.moduleRepo.GetByCourse(ctx, courseID)
}

func (uc *ModuleUseCase) GetModule(ctx context.Context, id string) (*entity.Module, error) {
	return uc.moduleRepo.GetByID(ctx, id)
}

func (uc *ModuleUseCase) UpdateModule(ctx context.Context, module *entity.Module, instructorID string) error {
	existingModule, err := uc.moduleRepo.GetByID(ctx, module.ID)
	if err != nil {
		return err
	}

	course, err := uc.courseRepo.GetByID(ctx, existingModule.CourseID)
	if err != nil {
		return err
	}

	if course.InstructorID != instructorID {
		return entity.ErrUnauthorized
	}

	return uc.moduleRepo.Update(ctx, module)
}

func (uc *ModuleUseCase) DeleteModule(ctx context.Context, id string, instructorID string) error {
	existing, err := uc.moduleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	course, err := uc.courseRepo.GetByID(ctx, existing.CourseID)
	if err != nil {
		return err
	}

	if course.InstructorID != instructorID {
		return entity.ErrUnauthorized
	}

	return uc.moduleRepo.Delete(ctx, id)
}
