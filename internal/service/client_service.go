package service

import (
	"github.com/rustingoff/admin_panel_rep/internal/model"
	"github.com/rustingoff/admin_panel_rep/internal/repository"
	"gopkg.in/go-playground/validator.v9"
)

type ClientService interface {
	CreateClient(client model.Client) error
	UpdateClient(client model.ClientUpdate, clientID uint) error
	DeleteClient(clientID uint) error

	GetAllClients() ([]model.Client, error)
	GetClient(clientID uint) (model.Client, error)
}

type clientService struct {
	repo      repository.ClientRepository
	validator *validator.Validate
}

func GetClientService(repo repository.ClientRepository, v *validator.Validate) ClientService {
	return &clientService{repo, v}
}

func (cs *clientService) CreateClient(client model.Client) error {
	err := cs.validator.Struct(client)
	if err != nil {
		return err
	}
	return cs.repo.CreateClient(client)
}

func (cs *clientService) UpdateClient(client model.ClientUpdate, clientID uint) error {
	err := cs.validator.Struct(client)
	if err != nil {
		return err
	}
	return cs.repo.UpdateClient(client, clientID)
}

func (cs *clientService) DeleteClient(clientID uint) error {
	return cs.repo.DeleteClient(clientID)
}

func (cs *clientService) GetAllClients() ([]model.Client, error) {
	return cs.repo.GetAllClients()
}

func (cs *clientService) GetClient(clientID uint) (model.Client, error) {
	return cs.repo.GetClient(clientID)
}
