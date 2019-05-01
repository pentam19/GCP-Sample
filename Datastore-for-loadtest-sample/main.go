package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

type ResponseJSON struct {
	UUID string
}

// Preparation
//  export GO111MODULE=on
//  go mod init gae-firestore
// Local
//  dev_appserver.py app.yaml
// Deploy
//  gcloud app deploy --project [projectid] -v testapiv001
func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "hello")
	})

	// curl -X POST https://[project].com/friends01/1
	r.POST("/friends01/:id", AddFriendEList)
	// curl https://[project].com/friends01/d05040b2-423d-4f91-a958-11f580e156ef
	r.GET("/friends01/:uuid", GetFriendEList)

	// curl -X POST https://[project].com/friends02/1
	r.POST("/friends02/:id", AddFriendSList)
	// curl https://[project].com/friends02/d05040b2-423d-4f91-a958-11f580e156ef
	r.GET("/friends02/:uuid", GetFriendSList)

	http.Handle("/", r)
	appengine.Main()
}

func AddFriendEList(c *gin.Context) {
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

	// FriendList
	var fList []FriendList
	var fkeys []*datastore.Key
	for i := 0; i < 10; i++ {
		f := FriendList{
			FriendID: "friend-" + nameid + "-" + strconv.Itoa(i),
		}
		fList = append(fList, f)
		fkey := datastore.NewKey(ctx, "FriendList", f.FriendID, 0, pkey)
		fkeys = append(fkeys, fkey)
	}
	_, err = datastore.PutMulti(ctx, fkeys, fList)
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
