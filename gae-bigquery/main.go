package main

import (
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
	"google.golang.org/appengine"

	"github.com/gin-gonic/gin"
)

// Preparation
//  export GO111MODULE=off <- off !!!
//  go mod init project
// Local
//  dev_appserver.py app.yaml
// Deploy
//  gcloud app deploy --project [projectid] -v test-app-v001
func main() {
	r := gin.New()
	r.GET("/", hello)
	r.GET("/get-bq-dataset", getbq)
	r.GET("/get-bq-pubdata", getpubdata)

	http.Handle("/", r)
	appengine.Main()
}

func hello(c *gin.Context) {
	c.String(200, "hello path")
}

func getbq(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	// Get the list of dataset names.
	names, err := datasets(ctx)
	if err != nil {
		c.String(500, fmt.Sprintf("err: %v", err))
	}

	c.JSON(200, names)
}

// datasets returns a list with the IDs of all the Big Query datasets visible
// with the given context.
func datasets(ctx context.Context) (ids []string, err error) {
	// Get the current application ID, which is the same as the project ID.
	projectID := appengine.AppID(ctx)

	// Create the BigQuery service.
	bq, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("could not create service: %v", err)
	}

	// Return a list of IDs.
	//var ids []string
	it := bq.Datasets(ctx)
	for {
		ds, err := it.Next()
		if err == iterator.Done {
			return ids, nil
		} else if err != nil {
			return nil, err
		}
		ids = append(ids, ds.DatasetID)
	}
}

func getpubdata(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	values, err := fetchBigQueryData(ctx)
	if err != nil {
		c.String(500, fmt.Sprintf("err: %v", err))
	}

	c.JSON(200, values)
}

func fetchBigQueryData(ctx context.Context) (results []string, err error) {

	projectID := appengine.AppID(ctx)

    // contextとprojectIDを元にBigQuery用のclientを生成
    client, err := bigquery.NewClient(ctx, projectID)

    if err != nil {
		return nil, fmt.Errorf("bq NewClient err: %v", err)
    }

	QUERY := "SELECT * FROM `bigquery-public-data.usa_names.usa_1910_2013` WHERE name = 'Mary' LIMIT 10"
    // 引数で渡した文字列を元にQueryを生成
    q := client.Query(QUERY)

    // 実行のためのqueryをサービスに送信してIteratorを通じて結果を返す
    // itはIterator
    it, err := q.Read(ctx)

    if err != nil {
		return nil, fmt.Errorf("bq Read err: %v", err)
    }

    for {
        // BigQueryの結果から、中身を格納するためのBigQuery.Valueのsliceを宣言
        // BigQuery.Valueはinterface{}型
        var values []bigquery.Value

        err := it.Next(&values)
        if err == iterator.Done {
			return results, nil
        }

        if err != nil {
			return nil, fmt.Errorf("iterator Next err: %v", err)
        }

        //fmt.Println(values)
		results = append(results, fmt.Sprintf("val: %v ", values))
    }
}
