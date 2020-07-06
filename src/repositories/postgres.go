package repositories

import (
	"database/sql/driver"
	"fmt"
	"log"

	"syncTool/src/pkg/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	portiereDB *sqlx.DB
	hermesDB   *sqlx.DB
)

type JSONField []byte

type Candidate struct {
	Host string
	User string
}

func InitPostgres() {
	portierePassword := config.GetConfig("connection.postgres.portierePassword")
	portiereDbName := config.GetConfig("connection.postgres.portiereDbName")
	portiereMasterHost := config.GetConfig("connection.postgres.portiereMasterHost")
	portiereMasterUser := config.GetConfig("connection.postgres.portiereMasterUser")
	// init portiere instance (master, write)
	portiereMaxOpen := config.GetInteger("connection.postgres.portiereMaxOpen")
	portiereMaxIdle := config.GetInteger("connection.postgres.portiereMaxIdle")

	hermesPassword := config.GetConfig("connection.postgres.hermesPassword")
	hermesDbName := config.GetConfig("connection.postgres.hermesDbName")
	hermesHost := config.GetConfig("connection.postgres.hermesHost")
	hermesUser := config.GetConfig("connection.postgres.hermesUser")
	// init hermes instance
	hermesMaxOpen := config.GetInteger("connection.postgres.hermesMaxOpen")
	hermesMaxIdle := config.GetInteger("connection.postgres.hermesMaxIdle")

	if portiereMasterInstance, err := initInstance(portiereMasterUser, portierePassword, portiereMasterHost, portiereDbName, portiereMaxOpen, portiereMaxIdle); err != nil {
		log.Println("Portiere DB connection fail started")
	} else {
		log.Println("Portiere DB connection started")
		portiereDB = portiereMasterInstance
	}

	if hermesInstance, err := initInstance(hermesUser, hermesPassword, hermesHost, hermesDbName, hermesMaxOpen, hermesMaxIdle); err != nil {
		log.Println("Hermes DB connection fail started")
	} else {
		log.Println("Hermes DB connection started")
		hermesDB = hermesInstance
	}
}

func Check() error {
	// ping write
	if err := portiereDB.Ping(); err != nil {
		return err
	}

	// ping read
	if err := hermesDB.Ping(); err != nil {
		return err
	}
	return nil
}

func (j JSONField) Value() (driver.Value, error) {
	if j != nil {
		return []byte(j), nil
	}
	return nil, nil
}

func (j *JSONField) Scan(value interface{}) error {
	if b, ok := value.([]byte); ok && string(b) != "{}" && string(b) != "" {
		*j = b
	}
	return nil
}

func getPortiereMasterInstance() *sqlx.DB {
	return portiereDB
}

func getHermesInstance() *sqlx.DB {
	return hermesDB
}

func initInstance(user, password, host, dbName string, maxOpen, maxIdle int) (result *sqlx.DB, err error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, host, dbName)
	if result, err = sqlx.Open("postgres", connString); err != nil {
		return nil, err
	} else {
		result.SetMaxOpenConns(maxOpen)
		result.SetMaxIdleConns(maxIdle)
	}
	return result, nil
}
