package main

import (
	"context"
	"encoding/csv"
	"io"
	"library-manager/dbimport"
	"library-manager/es"
	"library-manager/logger"
	"os"
)

var sContext = context.Background()

func main() {
	var les = es.MustCreateFromConfig()
	esInfo, err := les.GetClient().Info()
	if err != nil {
		logger.Error.Fatalln(err)
		return
	}
	logger.Info.Println(esInfo)
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
	var esMarc dbimport.ESBookMarc
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
		logger.Info.Printf("handle %d -> %v", index, row)
		esMarc.MARCRecNo = row[0]
		esMarc.CallNo = row[1]
		esMarc.Tittle = row[2]
		esMarc.Author = row[3]
		esMarc.Publisher = row[4]
		pubYear, err := dbimport.ParsePubYear(row[5])
		if err != nil {
			logger.Error.Printf("%d %v", index, err)
			break
		}
		esMarc.PubYear = pubYear
		esMarc.ISBN = row[6]
		logger.Info.Printf("%d -> %s", index, esMarc.ToJSON())
	}
	logger.Info.Printf("all finished")
}
