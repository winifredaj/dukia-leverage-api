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

	//Load environment variables from.env file
	err := godotenv.Load(".env")
	if err != nil {
        log.Println("Warning: No .env file found. Using system environment variables.")
    }
	
	
	//Build connection strings
	dbHost := os.Getenv("DB_HOST")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")
    dbPort := os.Getenv("DB_PORT")
    dbSSLMode := os.Getenv("DB_SSLMODE")


	//Check if any environmentvariable is missing or empty
	missingVars := []string{}

	if dbHost == "" {missingVars = append(missingVars, "DB_HOST")}
	if dbUser == "" {missingVars = append(missingVars, "DB_USER")} 
	if dbPassword == "" {missingVars = append(missingVars, "DB_PASSWORD")} 
	if dbName == "" {missingVars = append(missingVars, "DB_NAME")}
	if dbPort == "" {missingVars = append(missingVars, "DB_PORT")} 
	if dbSSLMode == "" {missingVars = append(missingVars, "DB_SSL_MODE")}

	if len(missingVars) > 0 {
        log.Fatalf("Error: Missing required database environment variables: %v", missingVars)
    }

	//Build DSN (Database Source Name)
	dsn:= fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
        dbHost, dbUser, dbPassword, dbName, dbPort, dbSSLMode,
	)

	log.Println("Connecting to database with DSN:", dsn)

	//Connect to PostgresSQL server with retries

	//for i := 0; i < 3; i++ {
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			//DB = database
			//break
		log.Fatal("Error connecting to the database:", err)
		}
		DB = database
		fmt.Println("Database connected successfully!")
		//log.Println("Retrying database connection... attempt", i+1)

	if DB == nil {
		log.Fatal("Error:Database connection could not be initialized.")
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

		-- Enseure golholding table has "quantity" column
		IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'gold_holdings') THEN
        	IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='gold_holdings' AND column_name='quantity') THEN
            	ALTER TABLE gold_holdings ADD COLUMN quantity DECIMAL DEFAULT 0 NOT NULL;
        	END IF;
   		 END IF;

		-- Ensure leverage_transactions table has 'net_disbursed' column 
    	IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'leverage_transactions') THEN
        	IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='leverage_transactions' AND column_name='net_disbursed') THEN
            	ALTER TABLE leverage_transactions ADD COLUMN net_disbursed DECIMAL DEFAULT 0 NOT NULL;
        	END IF;
    	END IF;

	END $$
		`).Error
	if err != nil {
        log.Fatalf("Error creating enum types or adding columns:%v", err)
    }
	fmt.Println("Enum type and columns created successfully!")

}
