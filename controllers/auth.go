package controllers

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/phpfor/lilac-go/helpers"
	"github.com/phpfor/lilac-go/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

//LoginGet handles GET /login route
func LoginGet(c *gin.Context) {
	h := helpers.DefaultH(c)
	h["Title"] = "Login form"
	h["Active"] = "login"
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()
	c.HTML(http.StatusOK, "auth/login", h)
}

//LoginPost handles POST /login route, authenticates user
func LoginPost(c *gin.Context) {
	session := sessions.Default(c)
	user := &models.User{}
	if err := c.Bind(user); err != nil {
		session.AddFlash("Please, fill out form correctly.")
		session.Save()
		c.Redirect(http.StatusFound, "/login")
		return
	}

	userData, _ := models.GetUserByEmail(user.Email)
	if userData.Email == "" {
		logrus.Errorf("Login error, IP: %s, Email: %s", c.ClientIP(), user.Email)
		session.AddFlash("Email or password incorrect1")
		session.Save()
		c.Redirect(http.StatusFound, "/login")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(user.Password)); err != nil {
		logrus.Errorf("Login error, IP: %s, Email: %s, Password: %s", c.ClientIP(), user.Email, user.Password)
		logrus.Errorf("userDB.Password: %s", userData.Password)
		session.AddFlash("Email or password incorrect2")
		session.Save()
		c.Redirect(http.StatusFound, "/login")
		return
	}

	session.Set("UserID", userData.Email)
	session.Save()
	c.Redirect(http.StatusFound, "/")
}

//SignUpGet handles GET /signup route
func SignUpGet(c *gin.Context) {
	h := helpers.DefaultH(c)
	h["Title"] = "Basic GIN web-site signup form"
	h["Active"] = "signup"
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()
	c.HTML(http.StatusOK, "auth/signup", h)
}

//SignUpPost handles POST /signup route, creates new user
func SignUpPost(c *gin.Context) {
	session := sessions.Default(c)
	user := &models.User{}
	if err := c.Bind(user); err != nil {
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/signup")
		return
	}

	userDB, _ := models.GetUserByEmail(user.Email)
	if userDB.Email != "" {
		session.AddFlash("User exists")
		session.Save()
		c.Redirect(http.StatusFound, "/signup")
		return
	}
	//create user
	err := user.HashPassword()
	if err != nil {
		session.AddFlash("Error whilst registering user.")
		session.Save()
		logrus.Errorf("Error whilst registering user: %v", err)
		c.Redirect(http.StatusFound, "/signup")
		return
	}

	if err := user.Insert(); err != nil {
		session.AddFlash("Error whilst registering user.")
		session.Save()
		logrus.Errorf("Error whilst registering user: %v", err)
		c.Redirect(http.StatusFound, "/signup")
		return
	}

	session.Set("UserID", user.Email)
	session.Save()
	c.Redirect(http.StatusFound, "/")
	return
}

//LogoutGet handles GET /logout route
func LogoutGet(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("UserID")
	session.Save()
	c.Redirect(http.StatusSeeOther, "/")
}
