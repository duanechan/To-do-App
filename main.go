package main

import (
	"encoding/gob"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(GenerateKey()))

func init() {
	store.Options.HttpOnly = true
	store.Options.Secure = true
	store.Options.SameSite = http.SameSiteStrictMode
	store.Options.MaxAge = 86400

	gob.Register(&User{})
}

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("views/*.html")
	router.Static("public", "./public")

	router.GET("/favicon.ico", func(c *gin.Context) { c.File("public/favicon.ico") })

	auth := router.Group("/auth")
	auth.GET("/login", loginHandler)
	auth.POST("/login", loginHandler)
	auth.POST("/logout", logoutHandler)

	protected := router.Group("/")
	protected.Use(authRequired())
	protected.GET("/", indexHandler)

	router.Run("[::]:8100")
}

func authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, _ := store.Get(c.Request, "session")
		sessionID, idExists := session.Values["session_id"].(string)
		authenticated, authExists := session.Values["authenticated"].(bool)

		if !authExists || !authenticated || !idExists || sessionID == "" {
			c.Redirect(http.StatusFound, "/auth/login")
			c.Abort()
			return
		}

		c.Next()
	}
}

func indexHandler(c *gin.Context) {
	session, _ := store.Get(c.Request, "session")
	user := session.Values["user"].(*User)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Data": map[string]any{
			"User": user,
		},
		"SessionID": session.ID,
	})
}

func loginHandler(c *gin.Context) {
	switch c.Request.Method {
	case "GET":
		c.HTML(http.StatusOK, "login.html", gin.H{})
	case "POST":
		username := c.PostForm("username")
		password := c.PostForm("password")

		if username == "" || password == "" {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{"Error": "Username and password required"})
			return
		}

		user, err := Login(username, password)
		if err != nil {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"Error": err})
			return
		}

		session, _ := store.Get(c.Request, "session")
		session.Values["session_id"] = uuid.New().String()
		session.Values["user"] = user
		session.Values["authenticated"] = true

		if err = session.Save(c.Request, c.Writer); err != nil {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{"Error": "Could not make request: " + err.Error()})
			return
		}

		c.Redirect(http.StatusFound, "/")
	}
}

func logoutHandler(c *gin.Context) {
	session, _ := store.Get(c.Request, "session")
	session.Values["user"] = nil
	session.Values["authenticated"] = false

	if err := session.Save(c.Request, c.Writer); err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"Error": "Could not make request: " + err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/auth/login")
}
