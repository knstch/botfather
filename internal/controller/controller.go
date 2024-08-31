package controller

import "github.com/gin-gonic/gin"

func InitRouter(h Handler) *gin.Engine {
	router := gin.Default()

	api := router.Group("/v1/api")

	bot := api.Group("/bot")
	bot.POST("", h.WelcomeHandler)

	return router
}
