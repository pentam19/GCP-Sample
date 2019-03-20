package appname

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// dev_appserver.py app.yaml
// gcloud app deploy --project [projectid] -v testapiv001
func init() {
	r := gin.New()
	r.GET("/hello", hello)
	r.GET("/", func(c *gin.Context) {
		c.String(200, "hello")
	})

	http.Handle("/", r)
}

func hello(c *gin.Context) {
	c.String(200, "hello path")
}
