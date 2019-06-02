package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uniplaces/carbon"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

type UserData struct {
	UUID string
	Name string
	Age  int
}

type FriendList struct {
	FriendID string
}

type UserDataFriendListSlice struct {
	UUID       string
	Name       string
	FriendList []string
}

type FriendListNoParentE struct {
	FriendID string
	UUID2    string
}

type GetFriendListNoParentEParam struct {
	UUID1 string
	UUID2 string
}

type ResponseJSON struct {
	UUID      string
	LastUUID2 string
}

type UserTime struct {
	UUID1 string
	UUID2 string `datastore:",noindex"`
	//Name  string `datastore:",noindex"`
	Score int   `datastore:",noindex"`
	Time  int64 `datastore:",noindex"`
}

// Preparation
//  export GO111MODULE=on
//  go mod init gae-firestore
// Local
//  dev_appserver.py app.yaml
// Deploy
//  gcloud app deploy --project [projectid] -v testapiv001
//  gcloud app deploy --project [projectID] -v [version] index.yaml
func main() {
	//gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "hello")
	})

	// curl -X POST https://[project].com/user/1
	r.POST("/user/:id", AddUser)
	// curl -X POST https://[project].com/friends01/d05040b2-423d-4f91-a958-11f580e156ef
	r.POST("/friends01/:uuid", AddFriendEList)
	// curl https://[project].com/friends01/d05040b2-423d-4f91-a958-11f580e156ef
	r.GET("/friends01/:uuid", GetFriendEList)

	// curl -X POST https://[project].com/friends02/1
	r.POST("/friends02/:id", AddFriendSList)
	// curl https://[project].com/friends02/d05040b2-423d-4f91-a958-11f580e156ef
	r.GET("/friends02/:uuid", GetFriendSList)

	// curl -X POST https://[project].com/friends03/d05040b2-423d-4f91-a958-11f580e156ef
	r.POST("/friends03/:uuid", AddFriendNoPEList)
	// curl https://[project].com/friends03/d05040b2-423d-4f91-a958-11f580e156ef
	r.GET("/friends03/:uuid", GetFriendNoPEList)
	// curl -X POST https://[project].com/friends03-2 -d '{"UUID1": "a958-11f580e156ef", "UUID2": "aaa"}'
	r.POST("/friends03-2", GetFriendNoPEList2)

	// curl -X POST https://[project].com/usertime/d05040b2-423d-4f91-a958-11f580e156ef
	r.POST("/usertime/:uuid", AddUserTimeList)
	// curl https://[project].com/usertime/d05040b2-423d-4f91-a958-11f580e156ef
	r.GET("/usertime1/:uuid", GetUserTimeList)
	// curl https://[project].com/usertime2/d05040b2-423d-4f91-a958-11f580e156ef
	r.GET("/usertime2/:uuid", GetUserTimeList2)

	http.Handle("/", r)
	appengine.Main()
}

// test1
func AddUser(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	// Parent
	nameid := c.Param("id")
	uuid := uuid.New().String()
	u := UserData{
		UUID: uuid,
		Name: "name" + nameid,
		Age:  20,
	}
	pkey := datastore.NewKey(ctx, "UserFriend01", u.UUID, 0, nil)
	_, err := datastore.Put(ctx, pkey, &u)
	if err != nil {
		c.String(500, "err")
	}

	//c.String(200, pkey.String())
	c.JSON(200, ResponseJSON{UUID: uuid})
}

func AddFriendEList(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	// Parent
	uuid := c.Param("uuid")
	pkey := datastore.NewKey(ctx, "UserFriend01", uuid, 0, nil)

	// FriendList
	var fList []FriendList
	var fkeys []*datastore.Key
	for i := 0; i < 10; i++ {
		f := FriendList{
			FriendID: "friend-" + uuid + "-" + strconv.Itoa(i),
		}
		fList = append(fList, f)
		fkey := datastore.NewKey(ctx, "FriendList", f.FriendID, 0, pkey)
		fkeys = append(fkeys, fkey)
	}
	_, err := datastore.PutMulti(ctx, fkeys, fList)
	if err != nil {
		c.String(500, err.Error())
	}

	//c.String(200, pkey.String())
	c.JSON(200, ResponseJSON{UUID: uuid})
}

func GetFriendEList(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	uuid := c.Param("uuid")
	pkey := datastore.NewKey(ctx, "UserFriend01", uuid, 0, nil)

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

// test2
func AddFriendSList(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	nameid := c.Param("id")
	uuid := uuid.New().String()
	// FriendList
	var flist []string
	for i := 0; i < 10; i++ {
		flist = append(flist, "friend-"+nameid+"-"+strconv.Itoa(i))
	}

	s := UserDataFriendListSlice{
		UUID:       uuid,
		Name:       nameid,
		FriendList: flist,
	}
	pkey := datastore.NewKey(ctx, "UserFriend02", s.UUID, 0, nil)
	_, err := datastore.Put(ctx, pkey, &s)
	if err != nil {
		c.String(500, err.Error())
	}

	//c.String(200, pkey.String())
	c.JSON(200, ResponseJSON{UUID: uuid})
}

func GetFriendSList(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	uuid := c.Param("uuid")
	key := datastore.NewKey(ctx, "UserFriend02", uuid, 0, nil)

	friends := &UserDataFriendListSlice{}
	err := datastore.Get(ctx, key, friends)
	if err != nil {
		c.String(500, err.Error())
	}

	c.JSON(200, friends.FriendList)
}

// test3
func AddFriendNoPEList(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	// ParentKey
	uuid1 := c.Param("uuid")
	pkey := datastore.NewKey(ctx, "UserNoEntity", uuid1, 0, nil)

	// FriendList
	var fList []FriendListNoParentE
	var fkeys []*datastore.Key
	var uuid2 string
	for i := 0; i < 10; i++ {
		uuid2 = uuid.New().String()
		f := FriendListNoParentE{
			FriendID: "friend-" + uuid1 + "-" + strconv.Itoa(i),
			UUID2:    uuid2,
		}
		fList = append(fList, f)
		fkey := datastore.NewKey(ctx, "FriendListNoPE", f.FriendID, 0, pkey)
		fkeys = append(fkeys, fkey)
	}
	_, err := datastore.PutMulti(ctx, fkeys, fList)
	if err != nil {
		c.String(500, err.Error())
	}

	//c.String(200, pkey.String())
	c.JSON(200, ResponseJSON{UUID: uuid1, LastUUID2: uuid2})
}

func GetFriendNoPEList(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	uuid := c.Param("uuid")
	pkey := datastore.NewKey(ctx, "UserNoEntity", uuid, 0, nil)

	q := datastore.NewQuery("FriendListNoPE").Ancestor(pkey).KeysOnly()

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

func GetFriendNoPEList2(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	param := GetFriendListNoParentEParam{}
	c.BindJSON(&param)
	pkey := datastore.NewKey(ctx, "UserNoEntity", param.UUID1, 0, nil)

	q := datastore.NewQuery("FriendListNoPE").Ancestor(pkey).Filter("UUID2 =", param.UUID2).KeysOnly()

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

// test4
func AddUserTimeList(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	uuid1 := c.Param("uuid")

	referenceTime, _ := carbon.CreateFromFormat(carbon.DefaultFormat, "2019-01-01 10:00:00", "UTC")

	kindName := "UserTime"
	var i int64
	var list []UserTime
	var keys []*datastore.Key
	for i = 0; i < 60; i++ {
		uuid2 := uuid.New().String()

		data := UserTime{
			UUID1: uuid1,
			UUID2: uuid2,
			//Name:  "name",
			//Score: 100,
			Time: referenceTime.Unix() + i,
		}
		list = append(list, data)
		key := datastore.NewIncompleteKey(ctx, kindName, nil)
		keys = append(keys, key)
	}

	_, err := datastore.PutMulti(ctx, keys, list)
	if err != nil {
		c.String(500, err.Error())
	}

	c.JSON(200, ResponseJSON{UUID: uuid1})
}

// パフォーマンスNG
func GetUserTimeList(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	uuid1 := c.Param("uuid")

	referenceTime, _ := carbon.CreateFromFormat(carbon.DefaultFormat, "2019-01-01 10:00:30", "UTC")

	kindName := "UserTime"
	var data []UserTime
	q := datastore.NewQuery(kindName).Filter("UUID1 =", uuid1).Filter("Time >", referenceTime.Unix())
	_, err := q.GetAll(ctx, &data)
	if err != nil {
		c.String(500, err.Error())
	}

	//c.JSON(200, data)
	c.String(200, fmt.Sprintf("uuid: %v, len: %v", uuid1, len(data)))
}

// OK
func GetUserTimeList2(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	uuid1 := c.Param("uuid")

	referenceTime, _ := carbon.CreateFromFormat(carbon.DefaultFormat, "2019-01-01 10:00:30", "UTC")

	kindName := "UserTime"
	var data []UserTime
	q := datastore.NewQuery(kindName).Filter("UUID1 =", uuid1)
	_, err := q.GetAll(ctx, &data)
	if err != nil {
		c.String(500, err.Error())
	}
	var upperRefData []UserTime
	for _, u := range data {
		if u.Time > referenceTime.Unix() {
			upperRefData = append(upperRefData, u)
		}
	}

	//c.JSON(200, data)
	c.String(200, fmt.Sprintf("uuid: %v, len: %v", uuid1, len(upperRefData)))
}
