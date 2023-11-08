package router

import "github.com/gin-gonic/gin"

type Router interface {
	Run(addr ...string) error
}

func NewRouter(s *server.Server) Router {
	g := gin.Default()
	g.Use(middlewares.AuthMiddleware(s.Repository))
	g.GET("/ping", s.PingHandler)
	g.POST("/api/user/register", s.SignUp)
	g.POST("/api/user/login", s.Login)
	g.POST("/api/user/orders", s.UploadOrderHandler)
	g.GET("/api/user/orders", s.GetOrders)
	g.GET("/api/user/balance", s.GetUserBalance)
	g.GET("/api/user/withdrawals", s.GetUserWithdrawals)
	g.POST("/api/user/balance/withdraw", s.WithdrawHandler)
	return g
}
