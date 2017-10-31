package controllers

import (
	"net/http"

	"github.com/phpfor/lilac-go/helpers"
	"github.com/gin-gonic/gin"
)

//AdminGet handles GET /admin/ route
func AdminGet(c *gin.Context) {
	h := helpers.DefaultH(c)
	h["Title"] = "Admin dashboard"
	c.HTML(http.StatusOK, "admin/dashboard/index", h)
}
