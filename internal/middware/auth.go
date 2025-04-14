package middware

import (
	"context"
	"net/http"

	"github.com/VicShved/shorturl/internal/app"
	"github.com/VicShved/shorturl/internal/logger"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

func setAuthCook(w http.ResponseWriter, userID *string) {

	token, _ := app.GetJWTTokenString(userID)
	http.SetCookie(w, &http.Cookie{
		Name:  app.AuthorizationCookName,
		Value: token,
	})
}

func parseTokenUserID(tokenStr string) (*jwt.Token, string, error) {
	claims := &app.CustClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(app.ServerConfig.SecretKey), nil
	})
	return token, (*claims).User, err
}

// auth middleware
func AuthMiddleware(next http.Handler) http.Handler {
	authFunc := func(w http.ResponseWriter, r *http.Request) {
		var userID string
		var token *jwt.Token
		cook, err := r.Cookie(app.AuthorizationCookName)
		//  если нет куки, то создаю новый
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
		logger.Log.Debug("User ", zap.String("ID", userID))
		// next handler
		ctx := context.WithValue(r.Context(), app.ContextUser, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(authFunc)
}
