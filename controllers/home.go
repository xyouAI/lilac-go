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
	var filter models.DaoFilter
	count,_ := filter.GetPostsCount()
	currentPage, _ := strconv.Atoi(c.DefaultQuery("p","1"))
	limit := 10
	Pagination := helpers.NewPaginator(c,limit,count)
	list, err := filter.GetPostsByPage(currentPage,limit)
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
