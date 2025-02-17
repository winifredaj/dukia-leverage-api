package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//Build connection strings
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	//Connect to PostgresSQL server
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	DB = database
	fmt.Println("Database connected successfully!")

	// Manually create enum type before migrations
	err = DB.Exec(`DO $$
	BEGIN
	    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'leverage_status') THEN
		CREATE TYPE leverage_status AS ENUM ('pending','approved','active','defaulted','liquidated');
		END IF;
		
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'margincall_status') THEN
		CREATE TYPE margincall_status AS ENUM ('pending','resolved','defaulted');
		END IF;

	END $$
		`).Error
	if err != nil {
        log.Fatalf("Error creating enum type:%v", err)
    }
	fmt.Println("Enum type created successfully!")

}
