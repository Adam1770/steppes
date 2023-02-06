package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func clientAuthMiddleware(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader != "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Next()
}

func (rtr *Router) htmlAuthMiddleware(c *gin.Context) {
	if c.Request.ParseForm() != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")
	_, err := rtr.Client.Login(username, password)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "forbidden.html", gin.H{
			"title": "Forbidden",
		})
		c.Abort()
	}
	c.Next()
}
