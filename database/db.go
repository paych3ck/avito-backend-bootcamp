package database

import (
	"avito-backend-bootcamp/migrations"
	"avito-backend-bootcamp/models"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	connectionStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	var err error
	DB, err = sql.Open("postgres", connectionStr)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("error connecting to the database: %v", err)
	}

	migrations.Migrate(DB)
	return nil
}

func InitTestDB() error {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	testDbName := os.Getenv("TEST_DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	connectionStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, testDbName)

	var err error
	DB, err = sql.Open("postgres", connectionStr)
	if err != nil {
		return fmt.Errorf("error opening test database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("error connecting to the test database: %v", err)
	}

	migrations.Migrate(DB)
	return nil
}

func ClearTestDB() error {
	if DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	tables := []string{"flats", "houses", "users"}
	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE %s RESTART IDENTITY CASCADE;", table)
		_, err := DB.Exec(query)
		if err != nil {
			return fmt.Errorf("error clearing table %s: %v", table, err)
		}
	}
	return nil
}

func CreateUser(user *models.User) error {
	if DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	query := "INSERT INTO users (email, password, user_type) VALUES ($1, $2, $3) RETURNING id"
	err := DB.QueryRow(query, user.Email, user.Password, user.UserType).Scan(&user.ID)

	if err != nil {
		log.Printf("Error creating user: %v\n", err)
		return err
	}

	return nil
}

func CreateHouse(house *models.House) error {
	if DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	query := "INSERT INTO houses (address, year, developer, created_at, update_at) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	err := DB.QueryRow(query, house.Address, house.Year, house.Developer, house.CreatedAt, house.UpdateAt).Scan(&house.Id)
	if err != nil {
		log.Printf("Error creating user: %v\n", err)
		return err
	}

	return nil
}

func CreateFlat(flat *models.Flat) error {
	if DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	query := "INSERT INTO flats (house_id, flat_number, price, rooms, status) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	err := DB.QueryRow(query, flat.HouseId, flat.FlatNumber, flat.Price, flat.Rooms, flat.Status).Scan(&flat.Id)
	if err != nil {
		log.Printf("Error creating flat: %v\n", err)
		return err
	}
	return nil
}

func UpdateHouse(houseId int32) error {
	if DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	query := "UPDATE houses SET update_at = $1 WHERE id = $2"
	_, err := DB.Exec(query, time.Now(), houseId)
	if err != nil {
		log.Printf("Error updating house: %v\n", err)
		return err
	}
	return nil
}

func UpdateFlatStatus(id int32, status string) (*models.Flat, error) {
	if DB == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	query := "UPDATE flats SET status = $1 WHERE id = $2 RETURNING id, house_id, price, rooms, status"
	row := DB.QueryRow(query, status, id)

	var flat models.Flat
	err := row.Scan(&flat.Id, &flat.HouseId, &flat.Price, &flat.Rooms, &flat.Status)
	if err != nil {
		log.Printf("Error updating flat status: %v\n", err)
		return nil, err
	}

	return &flat, nil
}

func GetFlatByID(id int32) (*models.Flat, error) {
	if DB == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	query := "SELECT id, house_id, flat_number, price, rooms, status FROM flats WHERE id = $1"
	row := DB.QueryRow(query, id)

	var flat models.Flat
	err := row.Scan(&flat.Id, &flat.HouseId, &flat.FlatNumber, &flat.Price, &flat.Rooms, &flat.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("flat not found")
		}
		return nil, err
	}

	return &flat, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	if DB == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	user := &models.User{}
	query := "SELECT id, email, password, user_type FROM users WHERE email = $1"
	row := DB.QueryRow(query, email)

	if err := row.Scan(&user.ID, &user.Email, &user.Password, &user.UserType); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}

func GetFlatsByHouseID(houseID int, status string) ([]models.Flat, error) {
	if DB == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	var rows *sql.Rows
	var err error

	if status == "all" {
		query := "SELECT id, house_id, flat_number, price, rooms, status FROM flats WHERE house_id = $1"
		rows, err = DB.Query(query, houseID)
	} else {
		query := "SELECT id, house_id, flat_number, price, rooms, status FROM flats WHERE house_id = $1 AND status = $2"
		rows, err = DB.Query(query, houseID, status)
	}

	if err != nil {
		log.Printf("Error fetching flats: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var flats []models.Flat
	for rows.Next() {
		var flat models.Flat
		err := rows.Scan(&flat.Id, &flat.HouseId, &flat.FlatNumber, &flat.Price, &flat.Rooms, &flat.Status)
		if err != nil {
			log.Printf("Error scanning flat: %v\n", err)
			return nil, err
		}
		flats = append(flats, flat)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error with rows: %v\n", err)
		return nil, err
	}

	return flats, nil
}
