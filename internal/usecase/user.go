package usecase

import (
	"github.com/example/go-clean-architecture/internal/entity"
	"github.com/example/go-clean-architecture/internal/repository"
	"github.com/example/go-clean-architecture/pkg/utils"
)

// UserUsecase defines the interface for user business logic
type UserUsecase interface {
	CreateUser(req entity.UserRequest) (*entity.UserResponse, error)
	GetUserByID(id uint) (*entity.UserResponse, error)
	GetUserByEmail(email string) (*entity.UserResponse, error)
	GetAllUsers() ([]entity.UserResponse, error)
	UpdateUser(id uint, req entity.UserRequest) (*entity.UserResponse, error)
	DeleteUser(id uint) error
}

// userUsecase implements UserUsecase interface
type userUsecase struct {
	userRepo repository.UserRepository
}

// NewUserUsecase creates a new user usecase
func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user
func (u *userUsecase) CreateUser(req entity.UserRequest) (*entity.UserResponse, error) {
	// Check if user already exists
	existingUser, _ := u.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, &EmailAlreadyExistsError{Email: req.Email}
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create new user entity
	user := &entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}

	// Save user to repository
	if err := u.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Return response
	response := &entity.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return response, nil
}

// GetUserByID retrieves a user by ID
func (u *userUsecase) GetUserByID(id uint) (*entity.UserResponse, error) {
	user, err := u.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	response := &entity.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return response, nil
}

// GetUserByEmail retrieves a user by email
func (u *userUsecase) GetUserByEmail(email string) (*entity.UserResponse, error) {
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	response := &entity.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return response, nil
}

// GetAllUsers retrieves all users
func (u *userUsecase) GetAllUsers() ([]entity.UserResponse, error) {
	users, err := u.userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var responses []entity.UserResponse
	for _, user := range users {
		response := entity.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
		responses = append(responses, response)
	}

	return responses, nil
}

// UpdateUser updates a user
func (u *userUsecase) UpdateUser(id uint, req entity.UserRequest) (*entity.UserResponse, error) {
	// Get existing user
	user, err := u.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update user fields
	user.Name = req.Name
	user.Email = req.Email

	// Hash the new password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	// Save updated user
	if err := u.userRepo.Update(user); err != nil {
		return nil, err
	}

	// Return response
	response := &entity.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return response, nil
}

// DeleteUser deletes a user by ID
func (u *userUsecase) DeleteUser(id uint) error {
	return u.userRepo.Delete(id)
}

// EmailAlreadyExistsError represents an error when email already exists
type EmailAlreadyExistsError struct {
	Email string
}

func (e *EmailAlreadyExistsError) Error() string {
	return "user with email " + e.Email + " already exists"
}
