package tripleCrown

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"goJoe/internal/service"
	"os"
)

func startDB() (*sql.DB, error) {
	db, err := sql.Open(os.Getenv("DRIVERNAME"), os.Getenv("DATASOURCENAME"))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func uploadUserData(u service.UserReg) error {
	// make a call to db to see if user is there

	//if user is there, just leave - no more work

	//if user is not there, query the db to add all the relevant data
	return nil
}
