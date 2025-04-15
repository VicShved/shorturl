package app

import (
	"github.com/golang-jwt/jwt/v4"
)

// type TypeUserID string

// func (t TypeUserID) String() string {
// 	return fmt.Sprint(string(t))
// }

type CustClaims struct {
	jwt.RegisteredClaims
	UserID string
}

var AuthorizationCookName = "Authorization"
var SigningMethod = jwt.SigningMethodHS512
var ContextUser = "User"

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
