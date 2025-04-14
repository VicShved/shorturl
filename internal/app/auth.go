package app

import (
	"github.com/golang-jwt/jwt/v4"
)

type CustClaims struct {
	jwt.RegisteredClaims
	UserID string
}

var AuthorizationCookName = "AuthorizationCook"
var SigningMethod = jwt.SigningMethodHS512
var ContextUserIDKey = "UserID"

func GetJWTTokenString(userID *string) (string, error) {
	claim := CustClaims{
		UserID: *userID,
	}
	token := jwt.NewWithClaims(SigningMethod, claim)
	tokenStr, err := token.SignedString([]byte(ServerConfig.SecretKey))
	if err != nil {
		return "", nil
	}
	return tokenStr, err
}
