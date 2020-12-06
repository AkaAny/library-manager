package main

import (
	"context"
	"encoding/csv"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"io"
	"library-manager/es"
	"library-manager/logger"
	"library-manager/model"
	"library-manager/rest"
	"os"
	"strings"
)

var sContext = context.Background()

func main() {
	rest.InitRestAPI()
}

func doImport(les *es.LibraryES) {
	//index代表表名，各个json字段对应列，id对应主键
	fs, err := os.Open("dbimport/marc.csv")
	defer func() {
		err := fs.Close()
		if err != nil {
			logger.Error.Fatalln(err)
			return
		}
	}()
	var csvReader = csv.NewReader(fs)
	//第一行作为列头
	columns, err := csvReader.Read()
	if err != nil {
		logger.Error.Fatalln(err)
		return
	}
	logger.Info.Printf("columns:%v", columns)
	var index int64 = 0
	var esMarc model.ESBookMarc
	for {
		row, err := csvReader.Read()
		if err != nil && err != io.EOF {
			logger.Error.Fatalf("can not read, err is %+v", err)
			break
		}
		if err == io.EOF {
			break
		}
		index++
		//logger.Info.Printf("handle %d -> %v", index, row)
		esMarc.MARCRecNo = row[0]
		esMarc.CallNo = row[1]
		esMarc.Title = row[2]
		esMarc.Author = row[3]
		esMarc.Publisher = row[4]
		//目前版本先跳过出版日期和ISBN
		//esMarc.ISBN = row[6]

		//logger.Info.Printf("%d -> %s", index, esMarc.ToJSON())
		req := esapi.IndexRequest{
			Index:      "marc",
			DocumentID: esMarc.MARCRecNo,
			Body:       strings.NewReader(esMarc.ToJSON()),
			Refresh:    "true",
		}
		resp, err := req.Do(sContext, les.GetClient())
		if err != nil {
			panic(err)
		}
		if resp.StatusCode != 201 && resp.StatusCode != 200 {
			logger.Error.Println(resp.String())
			return
		}
		err = resp.Body.Close() //回收连接
		if err != nil {
			panic(err)
		}
		logger.Info.Println(index)
	}
	logger.Info.Printf("all finished")
}
