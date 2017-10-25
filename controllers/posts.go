package controllers

import (
	"fmt"
	"net/http"

	//"html/template"

	"github.com/Sirupsen/logrus"
	"github.com/phpfor/lilac-go/helpers"
	"github.com/phpfor/lilac-go/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	//"github.com/russross/blackfriday"
	//"gopkg.in/guregu/null.v3"
	//"log"
	//"github.com/revel/modules/db/app"
	//"gopkg.in/mgo.v2/bson"
)

//PostGet handles GET /posts/:id route
func PostGet(c *gin.Context) {
	post, err := models.GetPostBySlug(c.Param("slug"))
	//fmt.Println("Hello %s", post)
	if err != nil {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	h := helpers.DefaultH(c)
	h["Title"] = post.Title

	//h["Description"] = template.HTML(string(blackfriday.MarkdownCommon([]byte(post.Description))))
	h["Post"] = post
	//h["Active"] = fmt.Sprintf("posts/%d", post.ID)
	c.HTML(http.StatusOK, "posts/show", h)
}

//PostIndex handles GET /admin/posts route
func PostIndex(c *gin.Context) {
	list, err := models.GetPosts()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	h := helpers.DefaultH(c)
	h["Title"] = "List of blog posts"
	h["List"] = list
	h["Active"] = "posts"
	c.HTML(http.StatusOK, "posts/index", h)
}

//PostNew handles GET /admin/new_post route
func PostNew(c *gin.Context) {
	tags, _ := models.GetTags()
	h := helpers.DefaultH(c)
	h["Title"] = "New post entry"
	h["Active"] = "posts"
	h["Tags"] = tags
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()

	c.HTML(http.StatusOK, "posts/form", h)
}

//PostCreate handles POST /admin/new_post route
func PostCreate(c *gin.Context) {
	post := &models.Post{}
	if err := c.Bind(post); err != nil {
		session := sessions.Default(c)
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/admin/new_post")
		return
	}

	//if user, exists := c.Get("User"); exists {
	//	post.UserID = null.IntFrom(user.(*models.User).ID)
	//}
	if err := post.Insert(); err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/posts")
}

//PostEdit handles GET /admin/posts/:id/edit route
func PostEdit(c *gin.Context) {
	post, _ := models.GetPost(c.Param("id"))
	if post.ID == "" {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	tags, _ := models.GetTags()
	categorys, _ := models.GetCategorys()
	h := helpers.DefaultH(c)
	h["Title"] = "Edit post entry"
	h["Active"] = "posts"
	h["Post"] = post
	h["Tags"] = tags
	h["Categorys"] = categorys
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()
	c.HTML(http.StatusOK, "posts/form", h)
}

//PostUpdate handles POST /admin/posts/:id/edit route
func PostUpdate(c *gin.Context) {
	post := &models.Post{}
	if err := c.Bind(post); err != nil {
		session := sessions.Default(c)
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/admin/posts/edit/%s", c.Param("id")))
		return
	}

	if err := post.Update(c.Param("id")); err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/posts")
}

//PostDelete handles POST /admin/posts/:id/delete route
func PostDelete(c *gin.Context) {
	post, _ := models.GetPost(c.Param("id"))
	if err := post.Delete(c.Param("id")); err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/posts")
}
