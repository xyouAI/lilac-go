package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/Sirupsen/logrus"
	"github.com/phpfor/lilac-go/helpers"
	"github.com/phpfor/lilac-go/models"
	"strconv"
)

//HomeGet handles GET / route
func HomeGet(c *gin.Context) {
	currentPage, _ := strconv.Atoi(c.DefaultQuery("p","1"))
	Pagination := helpers.NewPaginator(c,10,20)
	list, err := models.GetPostsByPage(currentPage,1)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	h := helpers.DefaultH(c)
	h["List"] = list
	h["Title"] = "Welcome to use Lilac-go Blog"
	h["Active"] = "index"
	h["Pagination"] = Pagination
	c.HTML(http.StatusOK, "index", h)
}
