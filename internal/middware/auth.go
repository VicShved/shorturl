package middware

import (
	"context"
	"net/http"

	"github.com/VicShved/shorturl/internal/app"
	"github.com/VicShved/shorturl/internal/logger"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

// type TypeUserID string

// func (t TypeUserID) String() string {
// 	return fmt.Sprint(string(t))
// }

type CustClaims struct {
	jwt.RegisteredClaims
	UserID string
}

type contextKey int

const (
	ContextUser contextKey = iota
)

var AuthorizationCookName = "Authorization"
var SigningMethod = jwt.SigningMethodHS512

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

func parseTokenUserID(tokenStr string) (*jwt.Token, string, error) {
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
			token, userID, _ = parseTokenUserID(cook.Value)
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
