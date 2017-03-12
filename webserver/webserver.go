package webserver

import (
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

type Webserver struct {
	*gin.Engine
}

func New() *Webserver {
	s := Webserver{
		Engine: gin.New(),
	}

	r := multitemplate.New()
	r.AddFromFilesFuncs("user:show", getFuncMap(), "templates/layout.html", "templates/views/user.html")
	s.HTMLRender = r

	s.GET("/assets/*path", staticAssetHandler)

	s.GET("/users/:id", userHandler)

	return &s
}
