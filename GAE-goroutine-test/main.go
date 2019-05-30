package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

// Preparation
//  export GO111MODULE=on
//  go mod init goroutine-test
// Local
//  dev_appserver.py app.yaml
// Deploy
//  gcloud app deploy --project [projectid] -v testapiv001
func main() {
	r := gin.New()
	r.GET("/hello", hello)

	http.Handle("/", r)
	appengine.Main()
}

func testGoroutine(ctx context.Context) {
	time.Sleep(5 * time.Second)
	log.Infof(ctx, "Test Goroutine!!!")
}

func hello(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	go testGoroutine(ctx)
	c.String(200, "Hello Response!\n")
}
