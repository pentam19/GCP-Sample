package main

import (
	"math/rand"
	"net/http"
	"unsafe"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type UserData struct {
	UserID string
	Name   string
	Age    int
}

type ChkOverSizeUser struct {
	UserID string `datastore:",noindex"`
	Name   string `datastore:",noindex"`
	Age    int    `datastore:",noindex"`
}

type Task struct {
	TaskNo   int64
	TaskName string
}

type Size struct {
	Size int
}

const rsLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = rsLetters[rand.Intn(len(rsLetters))]
	}
	return string(b)
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
	// curl -X POST https://[project].com/entities -d '{"UserID": "abcXXX01","Name":"TestName","Age":80}'
	r.POST("/entities", AddEntityGrp)
	// curl https://[project].com/entities/abcXXX01
	r.GET("/entities/:id", GetEntityGrp)
	r.POST("/sizeover", AddSizeOverEntitiy)
	r.GET("/sizeover", GetMaxSizeEntitiy)

	http.Handle("/", r)
	appengine.Main()
}

func hello(c *gin.Context) {
	c.String(200, "hello path")
}

func AddUser(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	// 2019/3月現在
	// GCPコンソールで新しいプロジェクトを作成する
	// Datastoreの画面を開く
	// Firestoreへのアップグレードボタンを押すとDatastoreモードのFirestoreになる
	// 以下の通り、従来のDatastoreへアクセスするコードがそのまま使える
	data := UserData{
		UserID: "123abc",
		Name:   "Test User Name",
		Age:    80,
	}
	key := datastore.NewIncompleteKey(ctx, "userdata", nil)
	_, err := datastore.Put(ctx, key, &data)
	if err != nil {
		c.String(500, "err")
	}
	c.String(200, "result")
}

func AddEntityGrp(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	// Parent
	/*
		u := UserData{
			UserID: "123abc",
			Name:   "Test User Name",
			Age:    80,
		}
	*/
	var u UserData
	c.BindJSON(&u)
	pkey := datastore.NewKey(ctx, "User", u.UserID, 0, nil)
	_, err := datastore.Put(ctx, pkey, &u)
	if err != nil {
		c.String(500, "err")
	}

	// child 1
	t1 := Task{
		TaskNo:   1,
		TaskName: "TaskName001",
	}
	t1key := datastore.NewKey(ctx, "Task", "", 0, pkey)
	_, err = datastore.Put(ctx, t1key, &t1)
	if err != nil {
		c.String(500, "err")
	}

	// child 2
	t2 := Task{
		TaskNo:   2,
		TaskName: "TaskName002",
	}
	t2key := datastore.NewKey(ctx, "Task", "", 0, pkey)
	_, err = datastore.Put(ctx, t2key, &t2)
	if err != nil {
		c.String(500, "err")
	}

	c.String(200, pkey.String())
}

func GetEntityGrp(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	userId := c.Param("id")
	pkey := datastore.NewKey(ctx, "User", userId, 0, nil)

	q := datastore.NewQuery("Task").Ancestor(pkey)

	var task []Task
	_, err := q.GetAll(ctx, &task)
	if err != nil {
		c.String(500, err.Error())
	}

	c.JSON(200, task)
}

func AddSizeOverEntitiy(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	var s Size
	c.BindJSON(&s)

	// 1MiB = 1,048,576
	// Datastore 制限
	// エンティティ1,048,572 バイト（1 MiB～4 バイト） // 1Mib - 4
	// プロパティ　1,048,487 バイト（1 MiB～89 バイト）// 1Mib - 89
	// {"Size": 1048487} is NG
	// total 1048472 + 5 + 8 = 1048485 OK
	data := ChkOverSizeUser{
		UserID: RandString(5), // POST {"Size": 1048472} is OK
		//UserID: RandString(10), // POST {"Size": 1048467} is OK
		Name: RandString(s.Size),
		Age:  30, // 8byte
	}
	key := datastore.NewIncompleteKey(ctx, "userdataOverSize", nil)
	_, err := datastore.Put(ctx, key, &data)
	if err != nil {
		c.String(500, err.Error())
		//unspecified noindex
		// API error 1 (datastore_v3: BAD_REQUEST): The value of property "Name" is longer than 1500 bytes.ok
		//specified noindex. property over size.
		// API error 1 (datastore_v3: BAD_REQUEST): The value of property "Name" is longer than 1048487 bytes.ok
		//specified noindex. property not over size. entity over size.
		// API error 1 (datastore_v3: BAD_REQUEST): entity is too bigok
	}
	log.Infof(ctx, "Size struct: %v", unsafe.Sizeof(data)) // 40
	log.Infof(ctx, "Size UserID: %v", len(data.UserID))
	log.Infof(ctx, "Size   Name: %v", len(data.Name))
	log.Infof(ctx, "Size    Age: %v", unsafe.Sizeof(data.Age))
	c.String(200, "ok")
}

func GetMaxSizeEntitiy(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	q := datastore.NewQuery("userdataOverSize")

	var u []ChkOverSizeUser
	_, err := q.GetAll(ctx, &u)
	if err != nil {
		c.String(500, err.Error())
	}

	c.JSON(200, u)
}
