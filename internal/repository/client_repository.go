package repository

import (
	"github.com/rustingoff/admin_panel_rep/internal/model"
	"gorm.io/gorm"
)

type ClientRepository interface {
	CreateClient(client model.Client) error
	UpdateClient(client model.ClientUpdate, clientID uint) error
	DeleteClient(clientID uint) error

	GetAllClients() ([]model.Client, error)
	GetClient(clientID uint) (model.Client, error)
}

type clientRepository struct {
	db *gorm.DB
}

func GetClientRepository(db *gorm.DB) ClientRepository {
	return &clientRepository{db}
}

func (cr *clientRepository) CreateClient(client model.Client) error {
	err := cr.db.Debug().Model(&model.Client{}).Create(&client).Error
	return err
}

func (cr *clientRepository) UpdateClient(client model.ClientUpdate, clientID uint) error {
	err := cr.db.Debug().Model(&model.Client{}).Where("id = ?", clientID).Updates(&client).Error
	return err
}

func (cr *clientRepository) DeleteClient(clientID uint) error {
	err := cr.db.Debug().Delete(&model.Client{}, clientID).Error
	return err
}

func (cr *clientRepository) GetAllClients() ([]model.Client, error) {
	var clients []model.Client

	res := cr.db.Debug().Limit(-1).Find(&clients)
	return clients, res.Error
}

func (cr *clientRepository) GetClient(clientID uint) (model.Client, error) {
	var client model.Client

	res := cr.db.Debug().Where("id = ?", clientID).Find(&client)
	return client, res.Error
}
