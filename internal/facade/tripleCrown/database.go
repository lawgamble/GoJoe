package tripleCrown

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"goJoe/internal/service"
	"os"
	"strconv"
)

type Service struct {
	db *sql.DB
}

func startDB() (*sql.DB, error) {
	db, err := sql.Open(os.Getenv("DRIVERNAME"), os.Getenv("DATASOURCENAME"))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (s *Service) UploadUserData(u service.UserReg) error {
	discordId, _ := strconv.Atoi(u.DiscordID)
	score, _ := strconv.Atoi(u.Score)
	if u.PavlovName == "" {
		u.PavlovName = "null"
	}
	if u.ScoreStatus == "" {
		u.ScoreStatus = "null"
	}
	query := fmt.Sprintf("INSERT IGNORE INTO users (discordId, oculusName, pcl, pavlovName, score, scoreStatus) VALUES (%d, '%v', '%v', '%v', %d, '%v')", discordId, u.OculusName, u.Pcl, u.PavlovName, score, u.ScoreStatus)
	_, err := s.db.Query(query)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetUserByDiscordIdDB(u service.UserReg) (*sql.Rows, error) {
	discordId, _ := strconv.Atoi(u.DiscordID)
	query := fmt.Sprintf("SLECET * FROM users WHERE discordId = %d", discordId)
	rows, err := s.db.Query(query)
	if err != nil {
		return rows, err
	}
	return rows, nil
}
