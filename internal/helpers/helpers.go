package helpers

import (
	"github.com/Ed-cred/SolarPal/internal/models"
	"golang.org/x/crypto/bcrypt"
)

//Compares the given slice of valid logins to the parsed data and returns 0 if the user does not exist.
func FindUser(list []models.User, compareUser *models.User) uint {
	for _, item := range list {
		if item.Username == compareUser.Username && bcrypt.CompareHashAndPassword([]byte(item.Password), []byte(compareUser.Password)) == nil {
			return item.ID
		}
	}
	return 0
}
