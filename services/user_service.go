package services

import (
	"asset-maintenance/models"
	"asset-maintenance/repositories"
	"errors"
)

type UserService interface {
	GetAllUsers() ([]models.User, error)
	GetUser(id string) (*models.User, error)
	UpdateUser(id string, input UpdateUserInput) (*models.User, error)
	DeleteUser(id string) error
}

type userService struct {
	userRepo repositories.UserRepository
}

type UpdateUserInput struct {
	Username string
	FullName string
	Role     string
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.FindAll()
}

func (s *userService) GetUser(id string) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *userService) UpdateUser(id string, input UpdateUserInput) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Validasi role
	if input.Role != "" {
		switch models.Role(input.Role) {
		case models.RoleEngineer, models.RoleLogistik, models.RoleManajer:
			user.Role = models.Role(input.Role)
		default:
			return nil, errors.New("invalid role value")
		}
	}

	// Validasi username
	if input.Username != "" {
		taken, err := s.userRepo.IsUsernameTaken(input.Username, id)
		if err != nil {
			return nil, err
		}
		if taken {
			return nil, errors.New("username already taken")
		}
		user.Username = input.Username
	}

	if input.FullName != "" {
		user.FullName = input.FullName
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) DeleteUser(id string) error {
	return s.userRepo.Delete(id)
}
