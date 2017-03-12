package webserver

import (
	"path"

	"github.com/gin-gonic/gin"
)

func staticAssetHandler(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, max-age=0, must-revalidate, no-store")
	c.File(path.Join("build", c.Param("path")))
}
