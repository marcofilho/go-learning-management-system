package usecase

import (
        "context"
        "github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
        "github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/repository"
        "github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/infrastructure/auth"
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

        user, err := entity.NewUser(email, password, firstName, lastName, role)
        if err != nil {
                return nil, err
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

        if err := user.CanAuthenticate(); err != nil {
                return "", nil, err
        }

        if err := user.ValidatePassword(password); err != nil {
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
