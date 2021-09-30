package db

import (
	"database/sql"
	"log"
)

// represents the phone_numbers table in the DB
type Phone struct {
	ID     int
	Number string
}

type DB struct {
	db *sql.DB
}

func Open(driverName, dataSource string) (*DB, error) {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (db *DB) Close() error {
	return db.db.Close()
}

func (db *DB) Seed() error {

	data := []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"1234567892",
		"(123)456-7892",
	}

	for _, number := range data {
		_, err := insertPhone(db.db, number)
		if err != nil {
			return err
		}
	}
	return nil

}

func insertPhone(db *sql.DB, phone string) (int, error) {
	statement := `INSERT INTO phone_number(value) VALUES($1) RETURNING id`
	var id int
	err := db.QueryRow(statement, phone).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func Reset(driverName, dataSource, dbName string) error {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}
	err = resetDB(db, dbName)
	if err != nil {
		return err
	}
	return db.Close()

}

func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)
	if err != nil {
		panic(err)
	}
	return nil
}

func resetDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		panic(err)
	}
	return createDB(db, name)
}

func Migrate(driverName, dataSource string) error {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}
	err = creatPhoneNumbersTable(db)
	if err != nil {
		return err
	}
	return db.Close()

}

func (db *DB) GetAllPhones() ([]Phone, error) {
	return getAllPhones(db.db)

}

func getAllPhones(db *sql.DB) ([]Phone, error) {
	rows, err := db.Query("SELECT id, value FROM phone_number")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allPhones []Phone
	for rows.Next() {
		var p Phone
		err := rows.Scan(&p.ID, &p.Number)
		if err != nil {
			log.Fatal(err)
		}
		allPhones = append(allPhones, p)
	}
	return allPhones, nil

}

func creatPhoneNumbersTable(db *sql.DB) error {
	statement := `CREATE TABLE IF NOT EXISTS phone_number (
		id SERIAL,
		value VARCHAR(255)
		)`
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) FindPhone(number string) (*Phone, error) {
	var p Phone
	statement := "SELECT id, value FROM phone_number WHERE value = $1"
	row := db.db.QueryRow(statement, number)
	err := row.Scan(&p.ID, &p.Number)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &p, nil
}

func (db *DB) UpdatePhone(p *Phone) error {
	statement := "UPDATE phone_number SET VALUE = $2 WHERE ID = $1"
	_, err := db.db.Exec(statement, p.ID, p.Number)
	return err
}

func (db *DB) DeletePhone(id int) error {
	statement := "DELETE FROM phone_number WHERE id = $1"
	_, err := db.db.Exec(statement, id)
	return err
}
