package router

import (
	"github.com/DavidGQK/go-loyalty-system/internal/middleware"
	"github.com/DavidGQK/go-loyalty-system/internal/server"
	"github.com/gin-gonic/gin"
)

type Router interface {
	Run(addr ...string) error
}

func NewRouter(s *server.Server) Router {
	g := gin.Default()
	g.GET("/ping", s.PingHandler)
	g.POST("/api/user/register", s.SignUp)
	g.POST("/api/user/login", s.Login)

	private := g.Group("/api/user")
	private.Use(middleware.AuthMiddleware(s.Repository))
	{
		private.GET("/orders", s.GetOrders)
		private.GET("/balance", s.GetUserBalance)
		private.GET("/withdrawals", s.GetUserWithdrawals)
		private.POST("/orders", s.UploadOrderHandler)
		private.POST("/balance/withdraw", s.WithdrawHandler)
	}
	return g
}
