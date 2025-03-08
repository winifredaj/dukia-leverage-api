package utils

import (
	"errors"
    "fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var userSecretKey = []byte(os.Getenv("USER_JWT_SECRET")) //Ensure this matches the middleware
var adminSecretKey = []byte(os.Getenv("ADMIN_JWT_SECRET")) //Ensure this matches the middleware

//GenerateToken generates a JWT token for a users and admin
func GenerateToken(userID uint, email string, role string)(string, error){
// Set expiration time for token
expirationTime := time.Now().Add(24 * time.Hour)
    claims := jwt.MapClaims{
        "user_id":  userID,
        "email":    email,
        "role": role, //Distinguish user vs admin
        "exp": expirationTime.Unix(), //Token valid for 24 hours
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    //Use admin secret for admin role
    if role == "admin" {
        return token.SignedString(adminSecretKey)
    }
    return token.SignedString(userSecretKey)
   
}

//ParseToken verifies the JWT token and extracts claims
func ValidateToken(tokenString string, secretKey []byte) (*jwt.Token, jwt.MapClaims, error){ 
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error){
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok{
        return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
    }
    return userSecretKey, nil //
    })


    if err != nil || !token.Valid {
        return nil, nil, errors.New("invalid token")
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return nil,nil, errors.New("invalid token claims")
    }
    return token, claims, nil

}

// HashPassword hashes a given password
func HashPassword(password string) (string, error) {
    bytes,err:= bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes),err
}

// CheckPasswordHash compares a hashed password with a plain-text password
func CheckPasswordHash(password, hashedPassword string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    return err == nil
}