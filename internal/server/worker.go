package server

import (
	"context"
	"fmt"
	"github.com/DavidGQK/go-loyalty-system/internal/logger"
	"github.com/DavidGQK/go-loyalty-system/internal/services/accrual"
	"github.com/DavidGQK/go-loyalty-system/internal/services/converter"
	"github.com/DavidGQK/go-loyalty-system/internal/store"
	"time"
)

const Delay = 3 * time.Second

func (s *Server) OrdersWorker() {
	var messages []*store.Order
	ticker := time.NewTicker(Delay)

	for {
		select {
		case msg := <-s.OrdersQueue:
			messages = append(messages, msg)
		case <-ticker.C:
			if len(messages) == 0 {
				continue
			}

			for _, order := range messages {
				res, err := accrual.CheckOrder(s.Config.AccrualHost, order.OrderNumber)
				logger.Log.Info(fmt.Sprintf("%v", res))
				if err != nil {
					logger.Log.Errorf("accrual check order error: %v", err)
					s.OrdersQueue <- order
					continue
				}

				if !hasFinalStatus(res.Status) {
					// система ещё не обработала заказ – отправляем обратно в очередь
					logger.Log.Debugw("accrual has not processed the order yet",
						"order_number", order.OrderNumber,
						"current_accrual_status", res.Status)
					s.OrdersQueue <- order
					continue
				}

				err = updateOrder(s, order, res.Status, res.Accrual)
				if err != nil {
					logger.Log.Error(err)
					s.OrdersQueue <- order
				}
			}
			messages = nil
		default:
			continue
		}
	}
}

func hasFinalStatus(status string) bool {
	return status == accrual.OrderStatusInvalid || status == accrual.OrderStatusProcessed
}

func updateOrder(s *Server, order *store.Order, accrualStatus string, accrualBonus float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	bonusAmount := converter.ConvertToCent(accrualBonus)

	return s.Repository.UpdateOrderStatus(ctx, order, accrualStatus, bonusAmount)
}