package server

import (
	"context"
	"github.com/DavidGQK/go-loyalty-system/internal/logger"
	"github.com/DavidGQK/go-loyalty-system/internal/services/converter"
	"github.com/DavidGQK/go-loyalty-system/internal/services/mycrypto"
	"github.com/DavidGQK/go-loyalty-system/internal/store"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (s *Server) PingHandler(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

func (s *Server) GetOrders(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	token := mycrypto.HashFunc(authHeader)
	user, err := s.Repository.FindUserByToken(ctx, token)
	if err != nil {
		logger.Log.Errorf("find user error: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	orders, err := s.Repository.FindOrdersByUserID(ctx, user.ID)
	if err != nil {
		logger.Log.Errorf("find orders error: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if len(orders) == 0 {
		c.JSON(http.StatusNoContent, gin.H{})
		return
	}

	var response GetOrdersResponse
	for _, order := range orders {
		orderResp := OrderResponse{
			Number:     order.OrderNumber,
			Status:     order.Status,
			UploadedAt: order.CreatedAt.Format(time.RFC3339),
		}

		ac := converter.ConvertFromCent(order.BonusAmount) // FIXME: Не знаю, как сделать по-другому
		if ac != 0 {
			orderResp.Accrual = ac
		}

		response = append(response, orderResp)
	}
	c.JSON(http.StatusOK, response)
}

func (s *Server) GetUserBalance(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	token := mycrypto.HashFunc(authHeader)
	user, err := s.Repository.FindUserByToken(ctx, token)
	if err != nil {
		logger.Log.Errorf("find user error: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	withdrawalSum, err := s.Repository.GetWithdrawalSumByUserID(ctx, user.ID)
	if err != nil {
		logger.Log.Errorf("find bonus_transactions error: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, GetUserBalanceResponse{
		Current:   converter.ConvertFromCent(user.Bonuses),  // в БД храним в копейках
		Withdrawn: converter.ConvertFromCent(withdrawalSum), // в БД храним в копейках
	})
}

func (s *Server) GetUserWithdrawals(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	token := mycrypto.HashFunc(authHeader)
	user, err := s.Repository.FindUserByToken(ctx, token)
	if err != nil {
		logger.Log.Errorf("find user error: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	bonusTransactions, err := s.Repository.FindBonusTransactionsByUserID(ctx, user.ID)
	if err != nil {
		logger.Log.Errorf("find bonus_transactions error: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var response GetUserWithdrawalsResponse
	for _, tr := range bonusTransactions {
		if tr.Type == store.WithdrawalType {
			response = append(response, WithdrawalsResponse{
				Order:       tr.OrderNumber,
				Sum:         converter.ConvertFromCent(tr.Amount),
				ProcessedAt: tr.CreatedAt.Format(time.RFC3339),
			})
		}
	}

	if len(response) == 0 {
		c.JSON(http.StatusNoContent, gin.H{})
		return
	}

	c.JSON(http.StatusOK, response)
}
