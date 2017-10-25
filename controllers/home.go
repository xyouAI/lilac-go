package controllers

import (
	"net/http"

	"github.com/phpfor/lilac-go/helpers"
	"github.com/gin-gonic/gin"
)

//HomeGet handles GET / route
func HomeGet(c *gin.Context) {
	h := helpers.DefaultH(c)
	h["Title"] = "Welcome to use Lilac-go Blog"
	h["Active"] = "home"
	c.HTML(http.StatusOK, "home/show", h)
}
