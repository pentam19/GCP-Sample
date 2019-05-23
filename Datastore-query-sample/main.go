package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type UserData struct {
	UUID  string
	UUID2 string
	Name  string
	Age   int
}

// Preparation
//  export GO111MODULE=on
//  go mod init project
// Local
//  dev_appserver.py app.yaml
// Deploy
//  gcloud app deploy --project [projectid] -v testapiv001
func main() {
	//gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "hello")
	})

	// curl -X POST https://[project].com/data
	r.POST("/data", AddData)
	// curl https://[project].com/uuid2/d05040b2-423d-4f91-a958-11f580e156ef
	r.GET("/uuid2/:uuid1", GetUUID2)
	// curl https://[project].com/exists/9b0aea2b-70fa-4ce3-87f0-bb1391edad49
	r.GET("/exists/:uuid2", ExistsUUID2)

	http.Handle("/", r)
	appengine.Main()
}

func AddData(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	uuid1 := uuid.New().String()
	uuid2 := uuid.New().String()

	data := UserData{
		UUID:  uuid1,
		UUID2: uuid2,
		Name:  "name",
		Age:   30,
	}
	pkey := datastore.NewKey(ctx, "UserData", data.UUID, 0, nil)
	_, err := datastore.Put(ctx, pkey, &data)
	if err != nil {
		c.String(500, err.Error())
	}

	//c.String(200, pkey.String())
	c.JSON(200, data)
}

func GetUUID2(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	uuid1 := c.Param("uuid1")
	key := datastore.NewKey(ctx, "UserData", uuid1, 0, nil)

	q := datastore.NewQuery("UserData").Filter("__key__ =", key).Project("UUID2").Limit(1)
	var uuid2list []struct{ UUID2 string }
	_, err := q.GetAll(ctx, &uuid2list)
	if err != nil {
		c.String(500, err.Error())
	}
	uuid2 := uuid2list[0].UUID2

	result := struct{ UUID2 string }{UUID2: uuid2}
	c.JSON(200, result)
}

func ExistsUUID2(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	uuid2 := c.Param("uuid2")

	q := datastore.NewQuery("UserData").Filter("UUID2 =", uuid2).KeysOnly().Limit(1)
	keys, err := q.GetAll(ctx, nil)
	if err != nil {
		c.String(500, err.Error())
	}
	e := false
	msg := "not exists"
	if len(keys) > 0 {
		e = true
		msg = "exists"
	}
	log.Infof(ctx, "exists: %v, %v", e, msg)

	result := struct{ Exists bool }{Exists: e}
	c.JSON(200, result)
}
