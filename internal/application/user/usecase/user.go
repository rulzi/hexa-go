package usecase

import (
	"context"
	"time"

	"github.com/rulzi/hexa-go/internal/application/user/dto"
	domainuser "github.com/rulzi/hexa-go/internal/domain/user"
)

// UserUseCase handles all user operations (create, get, list, update, delete, login)
type UserUseCase struct {
	userRepo            domainuser.Repository
	passwordHasher      domainuser.PasswordHasher
	notificationService domainuser.NotificationService
	tokenGen            domainuser.TokenGenerator
}

// NewUserUseCase creates a new UserUseCase
func NewUserUseCase(
	userRepo domainuser.Repository,
	passwordHasher domainuser.PasswordHasher,
	notificationService domainuser.NotificationService,
	tokenGen domainuser.TokenGenerator,
) *UserUseCase {
	return &UserUseCase{
		userRepo:            userRepo,
		passwordHasher:      passwordHasher,
		notificationService: notificationService,
		tokenGen:            tokenGen,
	}
}

func validateCreateUserRequest(req dto.CreateUserRequest) error {
	if req.Name == "" {
		return domainuser.NewNameRequired()
	}
	if req.Email == "" {
		return domainuser.NewEmailRequired()
	}
	if req.Password == "" {
		return domainuser.NewPasswordRequired()
	}
	if len(req.Password) < 6 {
		return domainuser.NewPasswordTooShort()
	}
	return nil
}

// Create creates a new user
func (uc *UserUseCase) Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	if err := validateCreateUserRequest(req); err != nil {
		return nil, err
	}

	existingUser, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, domainuser.NewEmailExists()
	}

	hashedPassword, err := uc.passwordHasher.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	newUser := &domainuser.User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := newUser.Validate(); err != nil {
		return nil, err
	}

	createdUser, err := uc.userRepo.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	_ = uc.notificationService.SendWelcomeEmail(ctx, createdUser.Email, createdUser.Name)

	return &dto.UserResponse{
		ID:        createdUser.ID,
		Name:      createdUser.Name,
		Email:     createdUser.Email,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
	}, nil
}

// Get retrieves a user by ID
func (uc *UserUseCase) Get(ctx context.Context, id int64) (*dto.UserResponse, error) {
	userEntity, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if userEntity == nil {
		return nil, domainuser.NewUserNotFound()
	}

	return &dto.UserResponse{
		ID:        userEntity.ID,
		Name:      userEntity.Name,
		Email:     userEntity.Email,
		CreatedAt: userEntity.CreatedAt,
		UpdatedAt: userEntity.UpdatedAt,
	}, nil
}

// List lists users with pagination
func (uc *UserUseCase) List(ctx context.Context, limit, offset int) (*dto.ListUsersResponse, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	users, err := uc.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := uc.userRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	userResponses := make([]dto.UserResponse, len(users))
	for i, u := range users {
		userResponses[i] = dto.UserResponse{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		}
	}

	return &dto.ListUsersResponse{
		Users:  userResponses,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

// Update updates a user
func (uc *UserUseCase) Update(ctx context.Context, id int64, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	existingUser, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if existingUser == nil {
		return nil, domainuser.NewUserNotFound()
	}

	if req.Email != existingUser.Email {
		emailUser, err := uc.userRepo.GetByEmail(ctx, req.Email)
		if err == nil && emailUser != nil {
			return nil, domainuser.NewEmailExists()
		}
	}

	existingUser.Name = req.Name
	existingUser.Email = req.Email
	existingUser.UpdatedAt = time.Now()

	if req.Password != "" {
		hashedPassword, err := uc.passwordHasher.Hash(req.Password)
		if err != nil {
			return nil, err
		}
		existingUser.Password = hashedPassword
	}

	if err := existingUser.Validate(); err != nil {
		return nil, err
	}

	updatedUser, err := uc.userRepo.Update(ctx, existingUser)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        updatedUser.ID,
		Name:      updatedUser.Name,
		Email:     updatedUser.Email,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
	}, nil
}

// Delete deletes a user
func (uc *UserUseCase) Delete(ctx context.Context, id int64) error {
	existingUser, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existingUser == nil {
		return domainuser.NewUserNotFound()
	}

	return uc.userRepo.Delete(ctx, id)
}

// Login authenticates user and returns token
func (uc *UserUseCase) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	userEntity, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, domainuser.NewInvalidCredentials()
	}

	if userEntity == nil {
		return nil, domainuser.NewInvalidCredentials()
	}

	if !uc.passwordHasher.Verify(userEntity.Password, req.Password) {
		return nil, domainuser.NewInvalidCredentials()
	}

	token, err := uc.tokenGen.Generate(userEntity.ID, userEntity.Email)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        userEntity.ID,
			Name:      userEntity.Name,
			Email:     userEntity.Email,
			CreatedAt: userEntity.CreatedAt,
			UpdatedAt: userEntity.UpdatedAt,
		},
	}, nil
}
