package mysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/wojciechpawlinow/usermanagement/internal/config"
)

type dbConfig struct {
	User            string
	Password        string
	Host            string
	Port            int
	DBName          string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

type Connections struct {
	Read  *sql.DB
	Write *sql.DB
}

// GetConnections initializes and returns database connections for read and write operations.
func GetConnections(c config.Provider) (*Connections, error) {
	readDB, err := initConn(dbConfig{
		User:            c.GetString("DB_READ_USER"),
		Password:        c.GetString("DB_READ_PASSWORD"),
		Host:            c.GetString("DB_READ_HOST"),
		Port:            c.GetInt("DB_READ_PORT"),
		DBName:          c.GetString("DB_READ_NAME"),
		MaxOpenConns:    c.GetInt("DB_READ_MAX_OPEN_CONN"),
		MaxIdleConns:    c.GetInt("DB_READ_MAX_IDLE_CONN"),
		ConnMaxLifetime: 10 * time.Minute,
		ConnMaxIdleTime: 5 * time.Minute,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize read DB: %w", err)
	}

	writeDB, err := initConn(dbConfig{
		User:            c.GetString("DB_WRITE_USER"),
		Password:        c.GetString("DB_WRITE_PASSWORD"),
		Host:            c.GetString("DB_WRITE_HOST"),
		Port:            c.GetInt("DB_WRITE_PORT"),
		DBName:          c.GetString("DB_WRITE_NAME"),
		MaxOpenConns:    c.GetInt("DB_WRITE_MAX_OPEN_CONN"),
		MaxIdleConns:    c.GetInt("DB_WRITE_MAX_IDLE_CONN"),
		ConnMaxLifetime: 10 * time.Minute,
		ConnMaxIdleTime: 5 * time.Minute,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize write DB: %w", err)
	}

	return &Connections{
		Read:  readDB,
		Write: writeDB,
	}, nil
}

func initConn(c dbConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local", c.User, c.Password, c.Host, c.Port, c.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB connection: %w", err)
	}

	db.SetMaxOpenConns(c.MaxOpenConns)
	db.SetMaxIdleConns(c.MaxIdleConns)
	db.SetConnMaxLifetime(c.ConnMaxLifetime)
	db.SetConnMaxIdleTime(c.ConnMaxIdleTime)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	return db, nil
}
