package repository

import (
	"github.com/rustingoff/admin_panel_rep/internal/model"
	"gorm.io/gorm"
)

const _UserRole = 0

type UserRepository interface {
	CreateUser(user model.User) error
	DeleteUser(userID uint) error
	GetUserByUsername(username string) (model.User, error)
	GetAllUsers() ([]model.User, error)
	GetUser(userID uint) (model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func GetUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (cr *userRepository) CreateUser(user model.User) error {
	user.Role = _UserRole

	err := cr.db.Debug().Model(&model.User{}).Create(&user).Error
	return err
}

func (cr *userRepository) DeleteUser(userID uint) error {
	err := cr.db.Debug().Delete(&model.User{}, userID).Error
	return err
}

func (cr *userRepository) GetAllUsers() ([]model.User, error) {
	var users []model.User

	res := cr.db.Debug().Limit(-1).Find(&users)
	return users, res.Error
}

func (cr *userRepository) GetUser(userID uint) (model.User, error) {
	var user model.User

	res := cr.db.Debug().Where("id = ?", userID).Find(&user)
	return user, res.Error
}

func (cr *userRepository) GetUserByUsername(username string) (model.User, error) {
	var user model.User

	res := cr.db.Debug().Where("username = ?", username).Find(&user)
	return user, res.Error
}
