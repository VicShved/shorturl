package app

import (
	"github.com/golang-jwt/jwt/v4"
)

type CustClaims struct {
	jwt.RegisteredClaims
	User string
}

var AuthorizationCookName = "AuthorizationCook"
var SigningMethod = jwt.SigningMethodHS512
var ContextUser = "User"

func GetJWTTokenString(userID *string) (string, error) {
	claim := CustClaims{
		User: *userID,
	}
	token := jwt.NewWithClaims(SigningMethod, claim)
	tokenStr, err := token.SignedString([]byte(ServerConfig.SecretKey))
	if err != nil {
		return "", nil
	}
	return tokenStr, err
}
