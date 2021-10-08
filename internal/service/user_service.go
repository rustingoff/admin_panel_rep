package service

import (
	"github.com/rustingoff/admin_panel_rep/internal/model"
	"github.com/rustingoff/admin_panel_rep/internal/repository"
	"gopkg.in/go-playground/validator.v9"
)

type UserService interface {
	CreateUser(client model.User) error
	DeleteUser(clientID uint) error

	GetAllUsers() ([]model.User, error)
	GetUser(clientID uint) (model.User, error)
}

type userService struct {
	repo      repository.UserRepository
	validator *validator.Validate
}

func GetUserService(repo repository.UserRepository, v *validator.Validate) UserService {
	return &userService{repo, v}
}

func (cs *userService) CreateUser(user model.User) error {
	err := cs.validator.Struct(user)
	if err != nil {
		return err
	}
	return cs.repo.CreateUser(user)
}

func (cs *userService) DeleteUser(userID uint) error {
	return cs.repo.DeleteUser(userID)
}

func (cs *userService) GetAllUsers() ([]model.User, error) {
	return cs.repo.GetAllUsers()
}

func (cs *userService) GetUser(userID uint) (model.User, error) {
	return cs.repo.GetUser(userID)
}
