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

type UserInfo struct {
	UserID   string
	FriendID string
	Name     string
}

type FriendList struct {
	FriendID string
}

type FriendListSlice struct {
	Id         string
	FriendList []string
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

type MapUser struct {
	FriendID string
	ClearNum int
}

type MultiID struct {
	UserID1 string
	UserID2 string
	Attr1   string
	Attr2   string
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

	r.GET("/", func(c *gin.Context) {
		c.String(200, "hello")
	})

	// curl -X POST https://[project].com/entities -d '{"UserID": "abcXXX01","Name":"TestName","Age":80}'
	r.POST("/entities", AddEntityGrp)
	// curl https://[project].com/entities/abcXXX01
	r.GET("/entities/:id", GetEntityGrp)
	r.POST("/sizeover", AddSizeOverEntitiy)
	r.GET("/sizeover", GetMaxSizeEntitiy)

	// curl -X POST https://[project].com/friends -d '{"UserID": "abcXXX01","Name":"TestName","FriendID": "aaaaaaa"}'
	r.POST("/friends", AddFriendList)
	// curl https://[project].com/friends/abcXXX01
	r.GET("/friends/:id", GetFriendList)

	// curl -X POST https://[project].com/friendlists -d '{"Id": "abcXXX01","FriendList":["f01", "f02"]}'
	r.POST("/friendlists", AddSlice)
	// curl -X PUT https://[project].com/friendlists -d 'f03'
	r.PUT("/friendlists/:id", AddSliceElm)
	// curl https://[project].com/friendlists/abcXXX01
	r.GET("/friendlists/:id", GetSlice)

	// curl -X POST https://[project].com/mapusers/001 -d '{"FriendID": "abcXXX01","ClearNum":100}'
	r.POST("/mapusers/:id", AddMapUser)
	// curl https://[project].com/mapusers/001
	r.GET("/mapusers/:id", GetMapUser)

	// curl -X POST https://[project].com/multiid -d '{"UserID1": "abcXXX01","UserID2": "abcXXX02","Attr1":"aaa","Attr2":"bbb"}'
	r.POST("/multiid", AddMultiID)
	// curl https://[project].com/multiid1/abcXXX01
	r.GET("/multiid1/:id", GetMultiID1)
	// curl https://[project].com/multiid2/abcXXX02
	r.GET("/multiid2/:id", GetMultiID2)
	// curl https://[project].com/multiid1only/abcXXX01
	r.GET("/multiid1only/:id", GetMultiID1IDOnly)

	http.Handle("/", r)
	appengine.Main()
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

func AddFriendList(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	// Parent
	/*
		u := UserInfo{
			UserID: "123abc",
			FriendID: "aaaaaaa",
			Name:   "Test User Name"
		}
	*/
	var u UserInfo
	c.BindJSON(&u)
	pkey := datastore.NewKey(ctx, "UserInfo", u.UserID, 0, nil)
	_, err := datastore.Put(ctx, pkey, &u)
	if err != nil {
		c.String(500, "err")
	}

	// FriendList 1
	f1 := FriendList{
		FriendID: "f00000001",
	}
	f1key := datastore.NewKey(ctx, "FriendList", f1.FriendID, 0, pkey)
	_, err = datastore.Put(ctx, f1key, &f1)
	if err != nil {
		c.String(500, "err")
	}

	// FriendList 1
	f2 := FriendList{
		FriendID: "f00000002",
	}
	f2key := datastore.NewKey(ctx, "FriendList", f2.FriendID, 0, pkey)
	_, err = datastore.Put(ctx, f2key, &f2)
	if err != nil {
		c.String(500, "err")
	}

	c.String(200, pkey.String())
}

func GetFriendList(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	userId := c.Param("id")
	pkey := datastore.NewKey(ctx, "UserInfo", userId, 0, nil)

	q := datastore.NewQuery("FriendList").KeysOnly().Ancestor(pkey)

	keys, err := q.GetAll(ctx, nil)
	if err != nil {
		c.String(500, err.Error())
	}

	keysSlice := make([]string, 0, len(keys))
	for _, key := range keys {
		keysSlice = append(keysSlice, key.StringID())
	}
	c.JSON(200, keysSlice)
}

func AddSlice(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	/*
		s := &FriendListSlice{
			Id:         "123abc",
			FriendList: []string{"f001", "f002"},
		}
	*/
	var s FriendListSlice
	c.BindJSON(&s)
	pkey := datastore.NewKey(ctx, "FriendListSlice", s.Id, 0, nil)
	_, err := datastore.Put(ctx, pkey, &s)
	if err != nil {
		c.String(500, err.Error())
	}

	c.String(200, pkey.String())
}

func AddSliceElm(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	/*
		s := &FriendListSlice{
			Id:         "123abc",
			FriendList: []string{"f001", "f002"},
		}
	*/
	userID := c.Param("id")
	buf := make([]byte, 2048)
	n, _ := c.Request.Body.Read(buf)
	fID := string(buf[0:n])
	key := datastore.NewKey(ctx, "FriendListSlice", userID, 0, nil)

	friends := &FriendListSlice{}
	err := datastore.Get(ctx, key, friends)
	if err != nil {
		c.String(500, err.Error())
	}
	friends.FriendList = append(friends.FriendList, fID)

	_, err = datastore.Put(ctx, key, friends)
	if err != nil {
		c.String(500, err.Error())
	}

	c.String(200, key.String())
}

func GetSlice(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	userID := c.Param("id")
	key := datastore.NewKey(ctx, "FriendListSlice", userID, 0, nil)

	friends := &FriendListSlice{}
	err := datastore.Get(ctx, key, friends)
	if err != nil {
		c.String(500, err.Error())
	}

	c.JSON(200, friends.FriendList)
}

func AddMapUser(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	/*
		s := &MapUser struct {
			FriendID: "f001"
			ClearNum: 200
		}
	*/
	mapId := c.Param("id")
	var m MapUser
	c.BindJSON(&m)
	pkey := datastore.NewKey(ctx, "Map"+mapId, m.FriendID, 0, nil)
	_, err := datastore.Put(ctx, pkey, &m)
	if err != nil {
		c.String(500, err.Error())
	}

	c.String(200, pkey.String())
}

func GetMapUser(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	mapId := c.Param("id")
	q := datastore.NewQuery("Map" + mapId).Limit(5).KeysOnly()
	keys, err := q.GetAll(ctx, nil)
	if err != nil {
		c.String(500, err.Error())
	}

	keysSlice := make([]string, 0, len(keys))
	for _, key := range keys {
		keysSlice = append(keysSlice, key.StringID())
	}
	c.JSON(200, keysSlice)
}

func AddMultiID(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	// Parent
	/*
		m := MultiID {
			UserID1: "123abc",
			UserID2: "abc123",
			attr1: "aaa",
			attr2: "bbb"
		}
	*/
	var m MultiID
	c.BindJSON(&m)
	key := datastore.NewKey(ctx, "MultiIDData", "", 0, nil)
	_, err := datastore.Put(ctx, key, &m)
	if err != nil {
		c.String(500, "err")
	}

	c.String(200, key.String())
}

func GetMultiID1(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	Id := c.Param("id")
	q := datastore.NewQuery("MultiIDData").Filter("UserID1 =", Id)

	var m []MultiID
	_, err := q.GetAll(ctx, &m)
	if err != nil {
		c.String(500, err.Error())
	}

	c.JSON(200, m)
}

func GetMultiID2(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	Id := c.Param("id")
	q := datastore.NewQuery("MultiIDData").Filter("UserID2 =", Id)

	var m []MultiID
	_, err := q.GetAll(ctx, &m)
	if err != nil {
		c.String(500, err.Error())
	}

	c.JSON(200, m)
}

func GetMultiID1IDOnly(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	Id := c.Param("id")
	// "UserID1", "Attr1"は射影できない
	q := datastore.NewQuery("MultiIDData").Filter("UserID2 =", Id).Project("UserID1", "Attr1")

	var ms []MultiID
	iter := q.Run(ctx)
	for {
		var m MultiID
		_, err := iter.Next(&m)
		if err == datastore.Done {
			break
		}
		if err != nil {
			c.String(500, err.Error())
			break
		}
		ms = append(ms, m)
	}

	c.JSON(200, ms)
}
