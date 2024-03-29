package database

import (
	"database/sql"
	"fmt"

	"github.com/leobowenwang/go_frame_proxy/internal/pkg/config"
	"github.com/leobowenwang/go_frame_proxy/internal/pkg/util"
	_ "github.com/lib/pq"
)

var DB *sql.DB

type IDatabase interface {
	GetUserdataFor(username string) (userdata Userdata, err error)
}

type DatabaseConn interface {
	connect(conf config.Config)
}

type Database struct {
	DB *sql.DB
}

type Userdata struct {
	HashedPassword []byte
	Roles          []byte
}

func Connect(config config.Config) error {
	// Initialize connection string.
	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=require connect_timeout=%d",
		config.Database.Host, config.Database.User, config.Database.Password, config.Database.Database, config.Database.Timeout)

	// Initialize connection object.
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	util.Logger.Debug("Successfully created connection to database")
	DB = db
	return nil
}

func (db *Database) GetUserdataFor(username string) (userdata Userdata, err error) {

	err = DB.QueryRow("SELECT password, roles from public.client WHERE username=$1", username).Scan(&userdata.HashedPassword, &userdata.Roles)

	if err != nil {
		if err == sql.ErrNoRows {
			return userdata, err
		}
		return userdata, err
	}

	return userdata, err
}
