package appname

import (
	"encoding/base64"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/pubsub/v1"
)

// dev_appserver.py app.yaml
// gcloud app deploy --project [project-id] -v testapiv001
func init() {
	r := gin.New()
	r.GET("/hello", hello)
	r.GET("/", func(c *gin.Context) {
		c.String(200, "hello")
	})
	r.POST("/publish", publish)

	http.Handle("/", r)
}

func hello(c *gin.Context) {
	c.String(200, "hello path")
}

func publish(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	hc, err := google.DefaultClient(ctx, pubsub.PubsubScope)
	if err != nil {
		c.String(500, "err: google.DefaultClient")
	}

	pubsubService, err := pubsub.New(hc)
	if err != nil {
		c.String(500, "err: pubsub.New")
	}

	t := time.Now()
	s := ""
	s = t.String()
	msgData := []byte("test!: date: " + s)
	resp, err := pubsubService.Projects.Topics.Publish(
		//"projects/YOUR-PROJECT-ID/topics/YOUR-TOPIC-ID",
		"projects/[project-id]/topics/test-topics",
		&pubsub.PublishRequest{
			Messages: []*pubsub.PubsubMessage{
				{
					Data: base64.StdEncoding.EncodeToString(msgData),
				},
			},
		},
	).Do()
	if err != nil {
		c.String(500, "err: pubsub.New")
	}
	c.String(200, "msg: "+resp.MessageIds[0])
	// pull
	// gcloud pubsub subscriptions pull --auto-ack projects/[project-id]/subscriptions/test-subscriptions
}
