package utils

import (
	"os"
	"time"
    "errors"

	"github.com/golang-jwt/jwt/v4"
)

var userSecretKey = []byte(os.Getenv("USER_JWT_SECRET")) //Ensure this matches the middleware
var adminSecretKey = []byte(os.Getenv("ADMIN_JWT_SECRET")) //Ensure this matches the middleware

func GenerateUserToken(UserID uint)(string, error){
	
    claims := &jwt.MapClaims{
        "user_id":  UserID,
        "exp": time.Now().Add(time.Hour * 24).Unix(), //Token valid for 24 hours
        "role": "user", //User role
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(userSecretKey)
}

func GenerateAdminToken(adminID uint)(string, error){
    // Retrieve the secret key from environment variables
    claims := &jwt.MapClaims{
        "admin_id": adminID,
        "exp": time.Now().Add(time.Hour * 24).Unix(), //Token valid for 24 hours
        "role": "admin", //Admin role
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(adminSecretKey)
}

//ParseToken verifies the JWT token and extracts claims
func ParseToken(tokenString string)(*jwt.MapClaims, error){ 
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC);!ok {
            return nil, errors.New("unexpected signing method")
        }
        return userSecretKey, nil  //Defaulting to user key
    })

    if err != nil || !token.Valid {
        return nil, errors.New("invalid token")
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return nil, errors.New("could not parse claims")
    }
    
    return &claims, nil

}