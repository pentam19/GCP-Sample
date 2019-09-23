package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	"google.golang.org/api/option"

	"cloud.google.com/go/datastore"
)

type Task struct {
	Desc string `datastore:"description"`
	Done bool   `datastore:"done"`
	id   int64  // The integer ID used in the datastore.
}

// Preparation
//  export GO111MODULE=on
//  go mod init project-name
// Local
//  dev_appserver.py app.yaml
// Deploy
//  gcloud app deploy --project [projectid] -v testapiv001
func main() {
	r := gin.New()

	r.POST("/task", AddTask)
	r.GET("/task", GetTask)

	http.Handle("/", r)
	appengine.Main()
}

func AddTask(c *gin.Context) {
	ctx := context.Background()

	projID := "[other projectid]"
	jsonPath := "./other-project.json"
	// https://cloud.google.com/docs/authentication/production?hl=ja
	client, err := datastore.NewClient(ctx, projID, option.WithCredentialsFile(jsonPath))
	if err != nil {
		log.Fatal(err)
		c.JSON(500, err)
	}

	task := &Task{
		Desc: "desc",
	}
	key := datastore.IncompleteKey("Task", nil)
	if _, err = client.Put(ctx, key, task); err != nil {
		c.JSON(500, err)
	}
	c.JSON(200, task)
}

func GetTask(c *gin.Context) {
	//ctx := appengine.NewContext(c.Request)
	ctx := context.Background()

	projID := "[other projectid]"
	jsonPath := "./other-project.json"
	// https://cloud.google.com/docs/authentication/production?hl=ja
	client, err := datastore.NewClient(ctx, projID, option.WithCredentialsFile(jsonPath))
	if err != nil {
		log.Fatalf("Could not create datastore client: %v", err)
		c.JSON(500, err)
	}

	var tasks []Task
	query := datastore.NewQuery("Task")
	keys, err := client.GetAll(ctx, query, &tasks)
	if err != nil {
		c.JSON(500, err)
	}

	for i, key := range keys {
		tasks[i].id = key.ID
	}
	c.JSON(200, tasks)
}
