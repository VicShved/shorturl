// middware
package middware

import (
	"context"
	"net/http"

	"github.com/VicShved/shorturl/internal/app"
	"github.com/VicShved/shorturl/internal/logger"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

// CustClaims struct
type CustClaims struct {
	jwt.RegisteredClaims
	UserID string
}

type contextKey int

// ContextUser
const (
	ContextUser contextKey = iota
)

// AuthorizationCookName
var AuthorizationCookName = "Authorization"

// SigningMethod
var SigningMethod = jwt.SigningMethodHS512

// GetJWTTokenString(userID *string)
func GetJWTTokenString(userID *string) (string, error) {
	claim := CustClaims{
		UserID: *userID,
	}
	token := jwt.NewWithClaims(SigningMethod, claim)
	tokenStr, err := token.SignedString([]byte(app.ServerConfig.SecretKey))
	if err != nil {
		return "", nil
	}
	return tokenStr, err
}

func setAuthCook(w http.ResponseWriter, userID *string) {

	token, _ := GetJWTTokenString(userID)
	http.SetCookie(w, &http.Cookie{
		Name:  AuthorizationCookName,
		Value: token,
	})
}

// ParseTokenUserID парсит jwt из строки
func ParseTokenUserID(tokenStr string) (*jwt.Token, string, error) {
	claims := &CustClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(app.ServerConfig.SecretKey), nil
	})
	return token, (*claims).UserID, err
}

// auth middleware
func AuthMiddleware(next http.Handler) http.Handler {
	authFunc := func(w http.ResponseWriter, r *http.Request) {
		var userID string
		var token *jwt.Token
		cook, err := r.Cookie(AuthorizationCookName)
		//  если нет куки, то создаю новую
		if err == http.ErrNoCookie {
			logger.Log.Debug("ErrNoCookie")
			userID, _ = app.GetNewUUID()
			setAuthCook(w, &userID)
		} else {
			token, userID, _ = ParseTokenUserID(cook.Value)
			// Если токен не валидный,  то создаю нвый userID
			if !token.Valid {
				logger.Log.Debug("Not valid token")
				userID, _ = app.GetNewUUID()
				setAuthCook(w, &userID)
			}
		}
		// Если кука не содержит ид пользователя, то возвращаю 401
		if userID == "" {
			logger.Log.Debug("Empty userID")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		logger.Log.Debug("User ", zap.String("ID", string(userID)))
		// добавляю userID в контекст
		ctx := context.WithValue(r.Context(), ContextUser, userID)
		// Вызываю след.обработчик
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(authFunc)
}
