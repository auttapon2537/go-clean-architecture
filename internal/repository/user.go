package repository

import (
	"github.com/example/go-clean-architecture/internal/entity"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(user *entity.User) error
	GetByID(id uint) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	GetAll() ([]entity.User, error)
	Update(user *entity.User) error
	Delete(id uint) error
}

// userRepository implements UserRepository interface
type userRepository struct {
	db Database
}

// NewUserRepository creates a new user repository
func NewUserRepository(db Database) UserRepository {
	return &userRepository{
		db: db,
	}
}

// Database interface for database operations
type Database interface {
	Create(value interface{}) error
	First(dest interface{}, conditions ...interface{}) error
	Find(dest interface{}, conditions ...interface{}) error
	Save(value interface{}) error
	Delete(value interface{}, conditions ...interface{}) error
}

// Create creates a new user
func (r *userRepository) Create(user *entity.User) error {
	return r.db.Create(user)
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.First(&user, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.First(&user, "email = ?", email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAll retrieves all users
func (r *userRepository) GetAll() ([]entity.User, error) {
	var users []entity.User
	err := r.db.Find(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// Update updates a user
func (r *userRepository) Update(user *entity.User) error {
	return r.db.Save(user)
}

// Delete deletes a user by ID
func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&entity.User{}, id)
}
