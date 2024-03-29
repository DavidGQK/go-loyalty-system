package middleware

import (
	"context"
	"fmt"
	"github.com/DavidGQK/go-loyalty-system/internal/logger"
	"github.com/DavidGQK/go-loyalty-system/internal/server"
	"github.com/DavidGQK/go-loyalty-system/internal/services/mycrypto"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func AuthMiddleware(r server.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if err := checkHeader(r, authHeader); err != nil {
			_ = c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		c.Next()
	}
}

func checkHeader(r server.Repository, header string) error {
	if len(header) == 0 {
		return fmt.Errorf("missing Authorization header")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	token := mycrypto.HashFunc(header)

	user, err := r.FindUserByToken(ctx, token)
	if err != nil {
		logger.Log.Errorf("find user error: %v", err)
		return fmt.Errorf("invalid token")
	}

	if user.TokenExpAt.Before(time.Now()) {
		return fmt.Errorf("token expired")
	}
	return nil
}
