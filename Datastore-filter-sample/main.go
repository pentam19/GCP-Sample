package main

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

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
	No    int
	Score int
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

	rand.Seed(time.Now().UnixNano())

	r.GET("/", func(c *gin.Context) {
		c.String(200, "hello")
	})

	// curl -X POST https://[project].com/data
	r.POST("/data", AddData)
	// curl https://[project].com/score-border-above/1
	r.GET("/score-border-above/:no", GetBorderAboveScores)
	// curl https://[project].com/score-border-below/1
	r.GET("/score-border-below/:no", GetBorderBelowScores)

	http.Handle("/", r)
	appengine.Main()
}

func AddData(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	for i := 0; i < 100; i++ {
		uuid1 := uuid.New().String()
		uuid2 := uuid.New().String()

		data := UserData{
			UUID:  uuid1,
			UUID2: uuid2,
			Name:  "name",
			No:    i / 10,
			Score: rand.Intn(10000),
		}
		log.Infof(ctx, "No   : %v", data.No)
		log.Infof(ctx, "Score: %v", data.Score)
		pkey := datastore.NewKey(ctx, "UserData", data.UUID, 0, nil)
		_, err := datastore.Put(ctx, pkey, &data)
		if err != nil {
			c.String(500, err.Error())
		}
	}

	c.String(200, "Add Complete")
}

func GetBorderAboveScores(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	var no int
	no, _ = strconv.Atoi(c.Param("no"))
	border := 5000

	var data []UserData
	q := datastore.NewQuery("UserData").Filter("No =", no).Filter("Score >", border)
	_, err := q.GetAll(ctx, &data)
	if err != nil {
		c.String(500, err.Error())
	}

	c.JSON(200, data)
}

func GetBorderBelowScores(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	var no int
	no, _ = strconv.Atoi(c.Param("no"))
	border := 5000

	var data []UserData
	q := datastore.NewQuery("UserData").Filter("No =", no).Filter("Score <", border)
	_, err := q.GetAll(ctx, &data)
	if err != nil {
		c.String(500, err.Error())
	}

	c.JSON(200, data)
}
