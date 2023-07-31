package repository

import "github.com/Ed-cred/SolarPal/internal/models"

type DBRepo interface {
	GetUsers() ([]models.User, error)
	CreateUser(user *models.User) error
	AddSolarArray(id uint, inputs models.RequiredInputs, opts models.OptionalInputs) (int, error)
	FetchSolarArrayData(userId uint, arrayId int) (models.RequiredInputs, models.OptionalInputs, error)
	UpdateSolarArrayData(arrayId int, userId uint, inputs *models.RequiredInputs, opts *models.OptionalInputs) error
	FetchUserArrays(userId uint) ([]int,error)
	RemoveSolarArrayData(userId uint, arrayId int) error
}
