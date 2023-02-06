package router

import (
	"steppes/client"

	"github.com/gin-gonic/gin"
)

type Router struct {
	Engine *gin.Engine
	Client client.CloudClient
}

func NewRouter(client client.CloudClient) *Router {
	router := &Router{
		Engine: gin.Default(),
		Client: client,
	}
	router.setupRouter()
	return router
}

func (router *Router) setupRouter() {
	router.Engine.Use(clientAuthMiddleware)
	router.setHtmlRoutes()
	router.setActionRoutes()
}
