package database

import (
	"context"
	"log"
	"time"

	"github.com/Ed-cred/SolarPal/internal/models"
)

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

func (m *SQLiteRepo) AddSolarArray(id uint, inputs models.RequiredInputs, opts models.OptionalInputs) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	statement := `INSERT INTO solar_array (azimuth, system_capacity, losses, array_type, module_type, tilt, address, user_id, 
		gcr, dc_ac_ratio, inv_eff, radius, dataset, soiling, albedo, bifaciality) 
		VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		returning array_id`
	var arrayId int
	err := m.DB.QueryRowContext(ctx, statement,
		inputs.Azimuth,
		inputs.SystemCapacity,
		inputs.Losses,
		inputs.ArrayType,
		inputs.ModuleType,
		inputs.Tilt,
		inputs.Address,
		id,
		opts.Gcr,
		opts.DcAcRatio,
		opts.InvEff,
		opts.Radius,
		opts.Dataset,
		opts.Soiling,
		opts.Albedo,
		opts.Bifaciality,
	).Scan(&arrayId)
	if err != nil {
		log.Println("Error inserting solar array parameters into database: ", err)
	}
	return arrayId, nil
}

func (m *SQLiteRepo) FetchSolarArrayData(userId uint, arrayId int) (models.RequiredInputs, models.OptionalInputs, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	var inputs models.RequiredInputs
	var opts models.OptionalInputs
	query := `SELECT azimuth, system_capacity, losses, array_type, module_type, tilt, address,
	gcr, dc_ac_ratio, inv_eff, radius, dataset, soiling, albedo, bifaciality, coalesce(lat,'') , coalesce(lon,'') FROM solar_array WHERE user_id = ? AND array_id = ?;`
	rows, err := m.DB.QueryContext(ctx, query, userId, arrayId)
	if err != nil {
		log.Println("Unable to retrieve array data:", err)
		return inputs, opts, err
	}
	for rows.Next() {

		err = rows.Scan(&inputs.Azimuth, &inputs.SystemCapacity, &inputs.Losses, &inputs.ArrayType, &inputs.ModuleType, &inputs.Tilt, &inputs.Address, &opts.Gcr, &opts.DcAcRatio, &opts.InvEff, &opts.Radius, &opts.Dataset, &opts.Soiling, &opts.Albedo, &opts.Bifaciality, &opts.Latitude, &opts.Longitude)
		if err != nil {
			return inputs, opts, err
		}
	}
	if err := rows.Err(); err != nil {
		return inputs, opts, err
	}
	return inputs, opts, nil
}

func (m *SQLiteRepo) FetchUserArrays(userId uint) ([]models.ArrayData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	var arraysData []models.ArrayData
	query := `SELECT array_id, address, coalesce(lat, ""), coalesce(lon, "") FROM solar_array WHERE user_id = ?`
	rows, err := m.DB.QueryContext(ctx, query, userId)
	if err != nil {
		log.Println("Error fetching array Ids: ", err)
		return arraysData, err
	}
	for rows.Next() {
		var arrayData models.ArrayData
		err = rows.Scan(&arrayData.ID, &arrayData.Adress, &arrayData.Latitude, &arrayData.Longitude)
		if err != nil {
			return arraysData, err
		}
		arraysData = append(arraysData, arrayData)
	}

	return arraysData, nil
}

func (m *SQLiteRepo) UpdateSolarArrayData(arrayId int, userId uint, inputs *models.RequiredInputs, opts *models.OptionalInputs) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	query := `UPDATE solar_array
	SET azimuth = ?, system_capacity = ?, losses = ?, array_type = ?, module_type = ?, tilt = ?, address = ?, 
	gcr = ?, dc_ac_ratio = ?, inv_eff = ?, radius = ?, dataset = ?, soiling = ?, albedo = ?, bifaciality = ?, lat=?, lon=?
	WHERE array_id = ? AND user_id = ?`
	_, err := m.DB.ExecContext(ctx, query,
		inputs.Azimuth,
		inputs.SystemCapacity,
		inputs.Losses,
		inputs.ArrayType,
		inputs.ModuleType,
		inputs.Tilt,
		inputs.Address,
		opts.Gcr,
		opts.DcAcRatio,
		opts.InvEff,
		opts.Radius,
		opts.Dataset,
		opts.Soiling,
		opts.Albedo,
		opts.Bifaciality,
		opts.Latitude,
		opts.Longitude,
		arrayId,
		userId,
	)
	if err != nil {
		log.Println("Error updating solar array data: ", err)
		return err
	}
	log.Println("Succesfully updated solar array data!")
	return nil
}

func (m *SQLiteRepo) RemoveSolarArrayData(userId uint, arrayId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	query := `DELETE FROM solar_array WHERE user_id = ? AND array_id = ?`
	_, err := m.DB.ExecContext(ctx, query, userId, arrayId)
	if err != nil {
		log.Println("Error removing solar array from database: ", err)
		return err
	}
	return nil
}
