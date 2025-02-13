package utils

import (
   "time"
   "github.com/golang-jwt/jwt/v4"
)

var secretKey = []byte("secretkey") //Ensure this matches the middleware

func GenerateToken(UserID uint)(string, error){
	
    claims := &jwt.MapClaims{
        "user_id":  UserID,
        "exp": time.Now().Add(time.Hour * 24).Unix(), //Token valid for 24 hourss
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(secretKey)
}