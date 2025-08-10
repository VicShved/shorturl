package gserver

import (
	"context"

	"github.com/VicShved/shorturl/internal/app"
	"github.com/VicShved/shorturl/internal/logger"
	"github.com/VicShved/shorturl/internal/middware"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func getNewUserIDToken() (string, string) {
	userID, _ := app.GetNewUUID()
	token, _ := middware.GetJWTTokenString(&userID)
	return userID, token
}

func authUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	logger.Log.Debug("In authUnaryInterceptor")
	var userID string
	var token *jwt.Token
	var tokenStr string
	var err error
	md, exists := metadata.FromIncomingContext(ctx)
	if !exists {
		logger.Log.Warn("authUnaryInterceptor hasnt metadata")
	}
	tokens := md.Get(middware.AuthorizationCookName)
	if len(tokens) == 0 {
		userID, tokenStr = getNewUserIDToken()
	}
	if len(tokens) > 0 {
		token, userID, err = middware.ParseTokenUserID(tokens[0])
		if err != nil {
			return nil, status.Errorf(codes.PermissionDenied, "Доступ запрещен")
		}
		// Если токен не валидный,  то создаю нвый userID
		if !token.Valid {
			logger.Log.Warn("Not valid token")
			userID, tokenStr = getNewUserIDToken()
		}
	}
	// Если кука не содержит ид пользователя, то возвращаю 401
	if userID == "" {
		logger.Log.Warn("Empty userID")
		return nil, status.Errorf(codes.PermissionDenied, "Доступ запрещен")
	}
	md.Set("userID", userID)
	md.Set(middware.AuthorizationCookName, tokenStr)
	newCtx := metadata.NewIncomingContext(ctx, md)

	logger.Log.Debug("Exit from authUnaryInterceptor", zap.String("token:", tokenStr))
	authHeader := metadata.Pairs(middware.AuthorizationCookName, tokenStr)
	_ = grpc.SetHeader(newCtx, authHeader)
	return handler(newCtx, req)

}
