package database

import (
	"log"

	"github.com/Ed-cred/SolarPal/internal/models"
)

// func (m *SQLiteRepo) GetRequiredInputs() {
// ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
// defer cancel()
// query := `SELECT azimuth, system_capacity, losses, array_type, module_type, tilt, adress
// 			FROM solar_array
// 			LEFT JOIN user ON solar_array.user_id = user.id`
// 	return
// }

func (m *SQLiteRepo) CreateUser(user *models.User) error {
	statement := `INSERT INTO user (username, password, email) VALUES (?, ?, ?)`
	_, err := m.DB.Exec(statement, user.Username, user.Password, user.Email)
	if err != nil {
		log.Println("Error creating user: ", err)
		return err
	}
	return nil

}

func (m *SQLiteRepo) GetUsers() ([]models.User, error) {
	query := `SELECT username, password, email FROM user`
	rows, err := m.DB.Query(query)
	var users []models.User
	if err != nil {
		return users, err
	}
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Username, &user.Password, &user.Email)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return users, err
	}

	return users, nil
}
