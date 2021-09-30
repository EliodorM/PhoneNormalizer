package main

import (
	"bytes"
	"fmt"

	phonedb "github.com/EliodorM/phone/db"

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
	errHandler(phonedb.Reset("postgres", psqlInfo, dbname)) // drop and recreate db

	psqlInfo = fmt.Sprintf("%s dbname =%s", psqlInfo, dbname)
	err := phonedb.Migrate("postgres", psqlInfo) // create tables
	errHandler(err)

	db, err := phonedb.Open("postgres", psqlInfo)
	errHandler(err)
	defer db.Close()

	errHandler(db.Seed())

	phones, err := db.GetAllPhones()
	errHandler(err)

	for _, p := range phones {
		fmt.Printf("Working on... %+v\n", p)
		number := normalize(p.Number)
		if number != p.Number {
			fmt.Println("Updating or removing...", number)
			existing, err := db.FindPhone(number)
			errHandler(err)

			if existing != nil {
				errHandler(db.DeletePhone(p.ID))
			} else {
				p.Number = number
				errHandler(db.UpdatePhone(&p))
			}
		} else {
			fmt.Println("No Changes required")
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

func normalize(phone string) string {
	var buf bytes.Buffer
	for _, ch := range phone {
		if ch >= '0' && ch <= '9' {
			buf.WriteRune(ch)
		}
	}
	return buf.String()
}
