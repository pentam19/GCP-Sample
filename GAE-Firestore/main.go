package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

type UserData struct {
	Name string
	Age  int
}

// Preparation
//  export GO111MODULE=on
//  go mod init gae-firestore
// Local
//  dev_appserver.py app.yaml
// Deploy
//  gcloud app deploy --project [projectid] -v testapiv001
func main() {
	r := gin.New()
	r.GET("/hello", hello)
	r.GET("/", func(c *gin.Context) {
		c.String(200, "hello")
	})
	r.POST("/user", AddUser)

	http.Handle("/", r)
	appengine.Main()
}

func hello(c *gin.Context) {
	c.String(200, "hello path")
}

func AddUser(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	// 2019/3月現在
	// GCPコンソールで新しいプロジェクトを作成し
	// Datastoreの画面を開く
	// Firestoreへのアップグレードボタンを押すとDatastoreモードのFirestoreになる
	// 以下の通り、従来のDatastoreへアクセスするコードがそのまま使える
	data := UserData{
		Name: "Test User Name",
		Age:  80,
	}
	key := datastore.NewIncompleteKey(c, "userdata", nil)
	_, err := datastore.Put(ctx, key, &data)
	if err != nil {
		c.String(500, "err")
	}
	c.String(200, "result")
}
