package router

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const templatesGlob = "static/*"

func (router *Router) setHtmlRoutes() {
	router.Engine.LoadHTMLGlob(templatesGlob)
	router.Engine.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "login.html", gin.H{
			"title": "Steppes - login",
		})
	})
	router.Engine.GET("/pipeline", func(ctx *gin.Context) {
		if ctx.Request.Referer() != "https://localhost:90/" {
			ctx.Redirect(http.StatusSeeOther, "/forbidden")
			return
		}
		data, err := router.Client.GetRuns()
		log.Println(data, err)
		ctx.HTML(http.StatusOK, "pipeline.html", data)
	})
	router.Engine.Use(router.htmlAuthMiddleware)
	router.Engine.POST("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusSeeOther, "/pipeline")
	})
}

func (router *Router) setActionRoutes() {

}
