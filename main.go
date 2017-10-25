//go:generate rice embed-go
package main

import (
	"flag"
	//"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/GeertJohan/go.rice"
	"github.com/Sirupsen/logrus"
	//"github.com/claudiu/gocron"
	"github.com/phpfor/lilac-go/controllers"
	"github.com/phpfor/lilac-go/helpers"
	"github.com/phpfor/lilac-go/models"
	"github.com/phpfor/lilac-go/system"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	//"github.com/utrack/gin-csrf"
)

//var db *mgo.Database

func main() {
	//gin.SetMode(gin.ReleaseMode)
	flag.Parse()

	setLogger()
	loadConfig()

	//Periodic tasks
	//gocron.Every(1).Day().Do(system.CreateXMLSitemap)
	//gocron.Start()

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()
	setTemplate(router) //initialize templates
	setSessions(router) //initialize session storage & use sessiom/csrf middlewares

	router.StaticFS("/public", http.Dir(system.PublicPath())) //better use nginx to serve assets (Cache-Control, Etag, fast gzip, etc)
	router.Use(SharedData())

	router.GET("/", controllers.HomeGet)
	router.NoRoute(controllers.NotFound)
	router.NoMethod(controllers.MethodNotAllowed)

	dao, err := models.NewDao()
	if err != nil {
		router.Use(controllers.ServiceError)
		return
	}
	defer dao.Close()

	if system.GetConfig().SignupEnabled {
		router.GET("/signup", controllers.SignUpGet)
		router.POST("/signup", controllers.SignUpPost)
	}
	router.GET("/login", controllers.LoginGet)
	router.POST("/login", controllers.LoginPost)
	router.GET("/logout", controllers.LogoutGet)

	router.GET("/posts/:slug", controllers.PostGet)
	router.GET("/category/:name", controllers.CategoryGet)
	router.GET("/tags/:name", controllers.TagGet)
	router.GET("/pages/:slug", controllers.PageGet)
	//router.GET("/archives/:year/:month", controllers.ArchiveGet)
	//router.GET("/rss", controllers.RssGet)
	//
	authorized := router.Group("/admin")
	authorized.Use(AuthRequired())
	{
		authorized.GET("/", controllers.AdminGet)
		authorized.POST("/upload", controllers.UploadPost) //image upload
		authorized.GET("/users", controllers.UserIndex)
		authorized.GET("/new_user", controllers.UserNew)
		authorized.POST("/new_user", controllers.UserCreate)
		authorized.GET("/users/edit/:id", controllers.UserEdit)
		authorized.POST("/users/edit/:id", controllers.UserUpdate)
		authorized.POST("/users/delete/:id", controllers.UserDelete)

		authorized.GET("/pages", controllers.PageIndex)
		authorized.GET("/new_page", controllers.PageNew)
		authorized.POST("/new_page", controllers.PageCreate)
		authorized.GET("/pages/edit/:slug", controllers.PageEdit)
		authorized.POST("/pages/edit/:slug", controllers.PageUpdate)
		authorized.POST("/pages/delete/:slug", controllers.PageDelete)

		authorized.GET("/posts", controllers.PostIndex)
		authorized.GET("/new_post", controllers.PostNew)
		authorized.POST("/new_post", controllers.PostCreate)
		authorized.GET("/posts/edit/:id", controllers.PostEdit)
		authorized.POST("/posts/edit/:id", controllers.PostUpdate)
		authorized.POST("/posts/delete/:id", controllers.PostDelete)

		authorized.GET("/category", controllers.CategoryIndex)
		authorized.GET("/new_category", controllers.CategoryNew)
		authorized.POST("/new_category", controllers.CategoryCreate)
		authorized.POST("/category/delete/:name", controllers.CategoryDelete)

		authorized.GET("/tags", controllers.TagIndex)
		authorized.GET("/new_tag", controllers.TagNew)
		authorized.POST("/new_tag", controllers.TagCreate)
		authorized.POST("/tags/delete/:name", controllers.TagDelete)
	}
	// Listen and server on 0.0.0.0:8080
	router.Run(":8080")
}

//setLogger initializes logrus logger with some defaults
func setLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stderr)
	if gin.Mode() == gin.DebugMode {
		logrus.SetLevel(logrus.InfoLevel)
	}
}

//setConfig loads config.json from rice box "config"
func loadConfig() {
	box := rice.MustFindBox("config")
	system.LoadConfig(box.MustBytes("config.json"))
}

//setTemplate loads templates from rice box "views"
func setTemplate(router *gin.Engine) {
	box := rice.MustFindBox("views")
	tmpl := template.New("").Funcs(template.FuncMap{
		"isActive":      helpers.IsActive,
		"stringInSlice": helpers.StringInSlice,
		"dateTime":      helpers.DateTime,
		"recentPosts":   helpers.RecentPosts,
		"tags":          helpers.Tags,
		"archives":      helpers.Archives,
	})

	fn := func(path string, f os.FileInfo, err error) error {
		if f.IsDir() != true && strings.HasSuffix(f.Name(), ".html") {
			var err error
			tmpl, err = tmpl.Parse(box.MustString(path))
			if err != nil {
				return err
			}
		}
		return nil
	}

	err := box.Walk("", fn)
	if err != nil {
		panic(err)
	}
	router.SetHTMLTemplate(tmpl)
}

//setSessions initializes sessions & csrf middlewares
func setSessions(router *gin.Engine) {
	config := system.GetConfig()
	//https://github.com/gin-gonic/contrib/tree/master/sessions
	store := sessions.NewCookieStore([]byte(config.SessionSecret))
	store.Options(sessions.Options{HttpOnly: true, MaxAge: 7 * 86400}) //Also set Secure: true if using SSL, you should though
	router.Use(sessions.Sessions("gin-session", store))
	//https://github.com/utrack/gin-csrf
	//router.Use(csrf.Middleware(csrf.Options{
	//	Secret: config.SessionSecret,
	//	ErrorFunc: func(c *gin.Context) {
	//		c.String(400, "CSRF token mismatch")
	//		c.Abort()
	//	},
	//}))
}

//+++++++++++++ middlewares +++++++++++++++++++++++

//SharedData fills in common data, such as user info, etc...
func SharedData() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		//ttt := session.Get("UserID")
		if UserID := session.Get("UserID"); UserID != nil {
			//logrus.Error("session get UserID : ", reflect.Type(UserID))
			user, _ := models.GetUserByEmail(UserID)
			if user.Email != "" {
				c.Set("User", user)
			}
		}
		if system.GetConfig().SignupEnabled {
			c.Set("SignupEnabled", true)
		}
		c.Next()
	}
}

//AuthRequired grants access to authenticated users, requires SharedData middleware
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if user, _ := c.Get("User"); user != nil {
			c.Next()
		} else {
			logrus.Warnf("User not authorized to visit %s", c.Request.RequestURI)
			c.HTML(http.StatusForbidden, "errors/403", nil)
			c.Abort()
		}
	}
}
