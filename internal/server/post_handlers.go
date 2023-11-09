package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/DavidGQK/go-loyalty-system/internal/logger"
	"github.com/DavidGQK/go-loyalty-system/internal/services/converter"
	"github.com/DavidGQK/go-loyalty-system/internal/services/mycrypto"
	"github.com/DavidGQK/go-loyalty-system/internal/services/validation"
	"github.com/DavidGQK/go-loyalty-system/internal/store"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

func (s *Server) UploadOrderHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	token := mycrypto.HashFunc(authHeader)
	user, err := s.Repository.FindUserByToken(ctx, token)
	if err != nil {
		logger.Log.Errorf("find user error: %v", err)
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	orderNumber, err := io.ReadAll(c.Request.Body)
	if err != nil || len(orderNumber) == 0 {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("request error"))
		return
	}
	if err := validation.LuhnValidate(string(orderNumber)); err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("invalid order_number"))
		return
	}

	order, err := s.Repository.FindOrderByOrderNumber(ctx, string(orderNumber))
	if err != nil {
		if err == store.ErrNowRows {
			order, err = s.Repository.CreateOrder(ctx, user.ID, string(orderNumber), store.OrderStatusNew)
			if err != nil {
				logger.Log.Errorf("create order error: %v", err)
				_ = c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("save order error %v", err))
				return
			}
			s.OrdersQueue <- order // кладём в очередь для фоновой обработки

			c.String(http.StatusAccepted, "order saved")
			return
		}
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	logger.Log.Infow("order already exists",
		"user_id", order.UserID,
		"order_number", order.OrderNumber)

	if order.UserID == user.ID {
		c.String(http.StatusOK, "order already exists")
		return
	}

	_ = c.AbortWithError(http.StatusConflict, fmt.Errorf("order already exists"))
}

func (s *Server) WithdrawHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	token := mycrypto.HashFunc(authHeader)
	user, err := s.Repository.FindUserByToken(ctx, token)
	if err != nil {
		logger.Log.Errorf("find user error: %v", err)
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var body WithdrawRequest
	decoder := json.NewDecoder(c.Request.Body)
	if err := decoder.Decode(&body); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if body.Sum > converter.ConvertFromCent(user.Bonuses) {
		_ = c.AbortWithError(http.StatusPaymentRequired, fmt.Errorf("insufficient funds"))
		return
	}
	if body.Sum == 0 {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid amount"))
		return
	}

	err = s.Repository.SaveWithdrawBonuses(ctx, user.ID, body.Order, converter.ConvertToCent(body.Sum))
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
