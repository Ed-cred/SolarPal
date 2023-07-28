package database

import (
	"context"
	"log"
	"time"

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
	query := `SELECT id, username, password, email FROM user`
	rows, err := m.DB.Query(query)
	var users []models.User
	if err != nil {
		return users, err
	}
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email)
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

// Returns created array id
func (m *SQLiteRepo) AddSolarArray(id uint, inputs models.RequiredInputs, opts models.OptionalInputs) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	statement := `INSERT INTO solar_array (azimuth, system_capacity, losses, array_type, module_type, tilt, adress, user_id, 
		gcr, dc_ac_ratio, inv_eff, radius, dataset, soiling, albedo, bifaciality) 
		VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.DB.ExecContext(ctx, statement,
		inputs.Azimuth,
		inputs.SystemCapacity,
		inputs.Losses,
		inputs.ArrayType,
		inputs.ModuleType,
		inputs.Tilt,
		inputs.Adress,
		id,
		opts.Gcr,
		opts.DcAcRatio,
		opts.InvEff,
		opts.Radius,
		opts.Dataset,
		opts.Soiling,
		opts.Albedo,
		opts.Bifaciality,
	)
	if err != nil {
		log.Println("Error inserting solar array parameters into database: ", err)
	}
	return nil
}

func (m *SQLiteRepo) FetchSolarArrayData(userId uint) ([]models.RequiredInputs, []models.OptionalInputs, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	var inputs []models.RequiredInputs
	var opts []models.OptionalInputs
	query := `SELECT azimuth, system_capacity, losses, array_type, module_type, tilt, adress,
	gcr, dc_ac_ratio, inv_eff, radius, dataset, soiling, albedo, bifaciality FROM solar_array WHERE user_id = ?`
	rows, err := m.DB.QueryContext(ctx, query, userId)
	if err != nil {
		log.Println("Unable to retrieve array data:", err)
		return inputs, opts, err
	}
	for rows.Next() {
		var input models.RequiredInputs
		var opt models.OptionalInputs

		err = rows.Scan(&input.Azimuth, &input.SystemCapacity, &input.Losses, &input.ArrayType, &input.ModuleType, &input.Tilt, &input.Adress, &opt.Gcr, &opt.DcAcRatio, &opt.InvEff, &opt.Radius, &opt.Dataset, &opt.Soiling, &opt.Albedo, &opt.Bifaciality)
		if err != nil {
			return inputs, opts, err
		}
		inputs = append(inputs, input)
		opts = append(opts, opt)
	}
	if err := rows.Err(); err != nil {
		return inputs, opts, err
	}
	return inputs, opts, nil
}
