package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/repository"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/infrastructure/auth"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserUseCase struct {
	userRepo      repository.UserRepository
	tokenProvider auth.TokenProvider
}

func NewUserUseCase(userRepo repository.UserRepository, tokenProvider auth.TokenProvider) *UserUseCase {
	return &UserUseCase{
		userRepo:      userRepo,
		tokenProvider: tokenProvider,
	}
}

func (uc *UserUseCase) Register(ctx context.Context, email, password, firstName, lastName string, role entity.UserRole) (*entity.User, error) {
	existing, _ := uc.userRepo.GetByEmail(ctx, email)
	if existing != nil {
		return nil, entity.ErrDuplicateEntry
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &entity.User{
		ID:           uuid.New().String(),
		Email:        email,
		PasswordHash: string(hashedPassword),
		FirstName:    firstName,
		LastName:     lastName,
		Role:         role,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (uc *UserUseCase) Login(ctx context.Context, email, password string) (string, *entity.User, error) {
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", nil, entity.ErrInvalidCredentials
	}
	if !user.IsActive {
		return "", nil, entity.ErrUnauthorized
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, entity.ErrInvalidCredentials
	}
	token, err := uc.tokenProvider.GenerateToken(user)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

func (uc *UserUseCase) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	return uc.userRepo.GetByID(ctx, id)
}

func (uc *UserUseCase) ListUsers(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	return uc.userRepo.List(ctx, limit, offset)
}

func (uc *UserUseCase) UpdateUser(ctx context.Context, user *entity.User) error {
	user.UpdatedAt = time.Now()
	return uc.userRepo.Update(ctx, user)
}

func (uc *UserUseCase) DeleteUser(ctx context.Context, id string) error {
	return uc.userRepo.Delete(ctx, id)
}
