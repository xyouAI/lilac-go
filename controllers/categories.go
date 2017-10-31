package controllers

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/phpfor/lilac-go/helpers"
	"github.com/phpfor/lilac-go/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Categories(c *gin.Context) {
	list, err := models.GetCategorys()
	if err != nil {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}

	h := helpers.DefaultH(c)
	h["Title"] = "Categories"
	h["Active"] = "categories"
	h["List"] = list
	h["Count"] = len(list)
	c.HTML(http.StatusOK, "categories/index", h)
}

//CategoryGet handles GET /Categorys/:name route
func CategoryGet(c *gin.Context) {
	Category, err := models.GetCategory(c.Param("name"))
	if err != nil {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	list, err := models.GetPostsByCategory(Category.Name)
	if err != nil {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}

	h := helpers.DefaultH(c)
	h["Title"] = Category.Name
	h["Category"] = Category
	h["Active"] = "categories"
	h["List"] = list
	c.HTML(http.StatusOK, "categories/show", h)
}

//CategoryIndex handles GET /admin/Categorys route
func CategoryIndex(c *gin.Context) {
	list, err := models.GetCategorys()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	h := helpers.DefaultH(c)
	h["Title"] = "List of Categorys"
	h["List"] = list
	h["Active"] = "categories"
	c.HTML(http.StatusOK, "admin/categories/index", h)
}

//CategoryNew handles GET /admin/new_Category route
func CategoryNew(c *gin.Context) {
	h := helpers.DefaultH(c)
	h["Title"] = "New Category"
	h["Active"] = "categories"
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()

	c.HTML(http.StatusOK, "admin/categories/form", h)
}

//CategoryCreate handles POST /admin/new_Category route
func CategoryCreate(c *gin.Context) {
	Category := &models.Category{}
	if err := c.Bind(Category); err != nil {
		session := sessions.Default(c)
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/admin/new_Category")
		return
	}

	if err := Category.Insert(); err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/categories")
}

//CategoryDelete handles POST /admin/Categorys/:name/delete route
func CategoryDelete(c *gin.Context) {
	Category, _ := models.GetCategory(c.Param("name"))
	if err := Category.Delete(); err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/categories")
}
