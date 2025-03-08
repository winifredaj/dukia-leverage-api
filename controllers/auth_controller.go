package controllers

import (
	"dukia-leverage-api/config"
	"dukia-leverage-api/models"
	"dukia-leverage-api/utils"
	"net/http"
	"os"
	"log"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser creates a new user
func RegisterUser(c *gin.Context) {
	var input struct {
		FirstName string `json:"first_name" binding:"required"`
		LastName  string `json:"last_name" binding:"required"`
		Email     string `json:"email" binding:"required,email"`
		Password  string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check if email already exists
	var existingUser models.User
	result := config.DB.Where("email =?", input.Email).First(&existingUser)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	// Hash the password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	user := models.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
        Email:    input.Email,
        Password: string(hashedPassword),
	}

	// Save user in database and return response
	result = config.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
		"user":    user,
    })
}

// User Login
func LoginUser(c *gin.Context) {
	var request  struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Find user by email
	var user models.User
	result := config.DB.Where("email =?", request.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	//Check password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	//Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email,"user")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

//AdminLogin handles login for admins and generates an admin JWT token
func LoginAdmin(c *gin.Context){
    var input struct {
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.BindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
        }
	
	
	//Retrieve admin credentials from environment variables
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	//Log the admin's input password for debugging purposes
	log.Printf("Input Password: %s", input.Password)
	log.Printf("Stored Hash: %s", adminPassword)

	if adminEmail == "" || adminPassword == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Admin credentials not set"})
		return
		
	}

    //Check admin credentials
	if input.Email != adminEmail || !utils.CheckPasswordHash(input.Password, adminPassword) {
        c.JSON(http.StatusUnauthorized, gin.H{"error":"Invalid email or password"})
        return
    }
        
    //Generate Admin Token 
    token, err := utils.GenerateToken(1, adminEmail, "admin")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to generate admin token"})
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "message": "Admin logged in successfully",
        "token": token,
    })


		// TEMPORARY DEBUG CODE
		inputPassword := "SecureAdmin123" // Replace with your actual password
		hash := "$2a$10$fsuT01/2RTlFCFEFoSunK.kw82MD2hBbV80hB12LdZq2vqN1iEgUW"
		match := utils.CheckPasswordHash(inputPassword, hash)
		log.Printf("DEBUG: Password Match? %v", match) // Should log "true"

}
