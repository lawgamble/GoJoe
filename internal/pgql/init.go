package pgql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var (
	host     = os.Getenv("HOST")
	port     = 5432
	user     = os.Getenv("USER")
	password = os.Getenv("PW")
)

func InitDatabase() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, "lvlzetxy")
	result, err := sql.Open("postgress", psqlInfo)
	if err != nil {
		log.Fatalf("DB ERROR : %s", err)
	}
	fmt.Println("Connected to the DB!")
	return result
}
