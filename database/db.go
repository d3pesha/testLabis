package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type ConfigDB struct {
	DBHost     string
	DBPort     string
	DBUsername string
	DBPassword string
	DBName     string
}

type Postgres struct {
	db *sqlx.DB
}

func NewPostgresConfig() *ConfigDB {
	err := godotenv.Load(".env") //"C:\\Users\\Reflex1on\\GolandProjects\\testLabis\\testLabis\\.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	return &ConfigDB{
		DBHost:     host,
		DBPort:     port,
		DBUsername: username,
		DBPassword: password,
		DBName:     dbname,
	}
}

func Connect(config *ConfigDB) Postgres {

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost,
		config.DBPort,
		config.DBUsername,
		config.DBPassword,
		config.DBName,
	)

	db, err := sqlx.Connect("pgx", connStr)

	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	return Postgres{db: db}
}
func (p *Postgres) GetDB() *sqlx.DB {
	return p.db
}
