package utils

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/ivinayakg/hypergroai_assignment/models"
)

func CreateJWT(user *models.User) (*string, error) {
	var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))
	var expiry = os.Getenv("JWT_EXPIRY")
	expiryTotal, err := strconv.Atoi(expiry)
	if err != nil {
		fmt.Println("Error:", err)
		expiryTotal = 21600
	}

	if expiryTotal <= 0 {
		return nil, fmt.Errorf("invalid expiry value: %s", expiry)
	}

	expirationTime := time.Now().Add(time.Duration(expiryTotal) * time.Second)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID,
		"email":  user.Email,
		"exp":    expirationTime.Unix(),
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

func VerifyJwt(tokenString string) (*map[string]string, error) {
	var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))
	if tokenString == "" {
		return nil, fmt.Errorf("invalid token value: %s", tokenString)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return &map[string]string{"userId": fmt.Sprint(claims["userId"]), "email": fmt.Sprint(claims["email"])}, nil
	} else {
		return nil, fmt.Errorf("failed to extract claims from token")
	}
}
