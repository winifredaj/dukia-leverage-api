package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"path/filepath"
)

var DB *gorm.DB

func ConnectDatabase() {

	rootPath, _ := filepath.Abs("..") // Moves up one directory
	envPath := filepath.Join(rootPath, ".env") // Path to.env file in the root directory

	//Load environment variables from.env file
	err := godotenv.Load(envPath)
	if err != nil {
        log.Fatal("Warning: No .env file found in root. Using system environment variables.")
    }
	
	
	//Build connection strings
	dbHost := os.Getenv("DB_HOST")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")
    dbPort := os.Getenv("DB_PORT")
    dbSSLMode := os.Getenv("DB_SSLMODE")

	if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" || dbPort == "" || dbSSLMode == "" {
        log.Fatal("Error: Missig required database environment variables.")
    }

	dsn:= fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
        dbHost, dbUser, dbPassword, dbName, dbPort, dbSSLMode,
	)
	log.Println("Connecting to database with DSN:", dsn)

	//Connect to PostgresSQL server
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	DB = database
	fmt.Println("Database connected successfully!")

	if DB == nil {
		log.Fatal("Error: config.DB is nil, database connection not initialized properly.")
    }

	// Manually create enum type before migrations
	err = DB.Exec(`DO $$
	BEGIN
	    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'leverage_status') THEN
		CREATE TYPE leverage_status AS ENUM ('pending','approved','active','defaulted','liquidated');
		END IF;
		
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'margincall_status') THEN
		CREATE TYPE margincall_status AS ENUM ('pending','resolved','defaulted');
		END IF;

		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'loan_status') THEN
		CREATE TYPE loan_status AS ENUM ('inactive','liquidated','resolved','defaulted');
		END IF;

	END $$
		`).Error
	if err != nil {
        log.Fatalf("Error creating enum type:%v", err)
    }
	fmt.Println("Enum type created successfully!")

}
