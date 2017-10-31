package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/phpfor/lilac-go/helpers"
)

//HomeGet handles GET / route
func AboutGet(c *gin.Context) {
	h := helpers.DefaultH(c)
	h["Title"] = "Welcome to use Lilac-go Blog"
	h["Active"] = "index"
	c.HTML(http.StatusOK, "about", h)
}
