package controllers

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/phpfor/lilac-go/helpers"
	"github.com/phpfor/lilac-go/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)


func Tags(c *gin.Context) {
	list, err := models.GetTags()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	h := helpers.DefaultH(c)
	h["Title"] = "List of tags"
	h["List"] = list
	h["Count"] = len(list)
	h["Active"] = "tags"
	c.HTML(http.StatusOK, "tags/index", h)
}

//TagGet handles GET /tags/:name route
func TagGet(c *gin.Context) {
	tag, err := models.GetTag(c.Param("name"))
	if err != nil {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	list, err := models.GetPostsByTag(tag.Name)
	if err != nil {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}

	h := helpers.DefaultH(c)
	h["Tag"] = tag
	h["Active"] = "tags"
	h["List"] = list
	c.HTML(http.StatusOK, "categories/show", h)
}

//TagIndex handles GET /admin/tags route
func TagIndex(c *gin.Context) {
	list, err := models.GetTags()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	h := helpers.DefaultH(c)
	h["Title"] = "List of tags"
	h["List"] = list
	h["Active"] = "tags"
	c.HTML(http.StatusOK, "admin/tags/index", h)
}

//TagNew handles GET /admin/new_tag route
func TagNew(c *gin.Context) {
	h := helpers.DefaultH(c)
	h["Title"] = "New tag"
	h["Active"] = "tags"
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()

	c.HTML(http.StatusOK, "admin/tags/form", h)
}

//TagCreate handles POST /admin/new_tag route
func TagCreate(c *gin.Context) {
	tag := &models.Tag{}
	if err := c.Bind(tag); err != nil {
		session := sessions.Default(c)
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/admin/new_tag")
		return
	}

	if err := tag.Insert(); err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/tags")
}

//TagDelete handles POST /admin/tags/:name/delete route
func TagDelete(c *gin.Context) {
	tag, _ := models.GetTag(c.Param("name"))
	if err := tag.Delete(); err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/tags")
}
