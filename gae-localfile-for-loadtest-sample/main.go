package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
	"google.golang.org/appengine"
)

type KeyValData struct {
	Key   string
	Value string
}

var keyValDataMap map[string]string = map[string]string{}

func init() {
	excelFileName := "./resource/SampleKeyValue.xlsx"
	/*
		Key	Value
		abc	abcValue
		bbb	bbbValue
	*/
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	for _, sheet := range xlFile.Sheets {
		//fmt.Println("SheetName: " + sheet.Name)
		for i, row := range sheet.Rows {
			if i == 0 {
				continue
			}
			var key string
			for j, cell := range row.Cells {
				text := cell.String()
				//fmt.Printf("%s\n", text)
				if j == 0 {
					key = text
					continue
				}
				keyValDataMap[key] = text
			}
		}
	}
}

// Preparation
//  export GO111MODULE=on
//  go mod init project
// Local
//  dev_appserver.py app.yaml
// Deploy
//  gcloud app deploy --project [projectid] -v test-app-v001
func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "hello")
	})

	r.GET("/values/:key", GetKeyValue)

	http.Handle("/", r)
	appengine.Main()
}

func GetKeyValue(c *gin.Context) {
	key := c.Param("key")

	if len(keyValDataMap[key]) == 0 {
		c.String(http.StatusBadRequest, "ERROR: Not found key!!!")
		return
	}

	c.JSON(200, KeyValData{Key: key, Value: keyValDataMap[key]})
}
