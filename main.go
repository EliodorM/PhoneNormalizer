package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "admin"
	dbname   = "postm"
)

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	db, err := sql.Open("postgres", psqlInfo) /// open a connection to the database SERVER
	// errHandler(err)

	// errHandler(err)
	//db.Close()
	err = resetDB(db, dbname)
	errHandler(err)
	db.Close()

	psqlInfo = fmt.Sprintf("%s dbname =%s", psqlInfo, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	errHandler(err)
	// Recreate (if exists) a database
	errHandler(err)

	err = creatPhoneNumbersTable(db)
	errHandler(err)

	_, err = insertPhone(db, "1234567890")
	errHandler(err)
	_, err = insertPhone(db, "123 456 7891")
	errHandler(err)
	_, err = insertPhone(db, "(123) 456 7892")
	errHandler(err)
	_, err = insertPhone(db, "(123) 456-7893")
	errHandler(err)
	_, err = insertPhone(db, "123-456-7894")
	errHandler(err)
	_, err = insertPhone(db, "123-456-7890")
	errHandler(err)
	_, err = insertPhone(db, "1234567892")
	errHandler(err)
	_, err = insertPhone(db, "(123)456-7892")
	errHandler(err)

	phones, err := getAllPhones(db)
	errHandler(err)

	for _, p := range phones {
		number := normalize(p.number)
		if number != p.number {
			existing, err := findPhone(db, number)
			errHandler(err)
			if existing != nil {
				errHandler(deletePhone(db, p.id))
			} else {
				p.number = number
				errHandler(updatePhone(db, p))
			}
		}
		fmt.Printf("%+v\n", p)
	}

	defer db.Close() // defer - action gets executed after all the functions around got executed
}

func errHandler(err error) { //handle general err
	if err != nil {
		panic(err)
	}
}

func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)
	errHandler(err)
	return nil
}

func resetDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	errHandler(err)
	return createDB(db, name)
}

func creatPhoneNumbersTable(db *sql.DB) error {
	statement := `CREATE TABLE IF NOT EXISTS phone_number (
		id SERIAL,
		value VARCHAR(255)
		)`
	_, err := db.Exec(statement)
	fmt.Println(err)
	errHandler(err)
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

type phone struct {
	id     int
	number string
}

func findPhone(db *sql.DB, number string) (*phone, error) {
	var p phone
	statement := "SELECT id, value FROM phone_number WHERE value = $1"
	row := db.QueryRow(statement, number)
	err := row.Scan(&p.id, &p.number)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &p, nil
}

func getAllPhones(db *sql.DB) ([]phone, error) {
	rows, err := db.Query("SELECT id, value FROM phone_number")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allPhones []phone
	for rows.Next() {
		var p phone
		err := rows.Scan(&p.id, &p.number)
		if err != nil {
			log.Fatal(err)
		}
		allPhones = append(allPhones, p)
	}
	return allPhones, nil

}
func updatePhone(db *sql.DB, p phone) error {
	statement := "UPDATE phone_number SET VALUE = $2 WHERE ID = $1"
	_, err := db.Exec(statement, p.id, p.number)
	return err
}

func deletePhone(db *sql.DB, id int) error {
	statement := "DELETE FROM phone_number WHERE id = $1"
	_, err := db.Exec(statement, id)
	return err
}

func normalize(phone string) string {
	var buf bytes.Buffer
	for _, ch := range phone {
		if ch >= '0' && ch <= '9' {
			buf.WriteRune(ch)
		}
	}
	return buf.String()
}
