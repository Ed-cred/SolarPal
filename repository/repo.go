package repository

import "github.com/Ed-cred/SolarPal/internal/models"

type DBRepo interface {
	GetUsers() ([]models.User, error)
	CreateUser(user *models.User) error
	AddSolarArray(id uint, inputs models.RequiredInputs, opts models.OptionalInputs) (int, error)
	FetchSolarArrayData(userId uint) ([]models.RequiredInputs, []models.OptionalInputs, error)
}
