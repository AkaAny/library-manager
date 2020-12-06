package dbimport

import (
	"bufio"
	"encoding/csv"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/pkg/errors"
	"io"
	"library-manager/logger"
	"library-manager/model"
	"mime/multipart"
)

type CSVToESHandler struct {
	bulkSession *BulkSession
}

func CreateCSVToESHandler(client *elasticsearch.Client) (*CSVToESHandler, error) {
	bulkSession, err := Create(client, "marc")
	if err != nil {
		return nil, err
	}
	return &CSVToESHandler{bulkSession: bulkSession}, nil
}

// 处理CSV导入请求，外部没有进行并发保护
func (handler CSVToESHandler) HandleCSVImport(csvStream *multipart.FileHeader) error {
	csvFile, err := csvStream.Open()
	if err != nil {
		logger.Error.Println(err)
		return err
	}
	defer csvFile.Close()
	var reader = csv.NewReader(bufio.NewReader(csvFile))
	index, err := handleRows(reader, handler.addToES)
	if err != nil {
		return errors.Wrapf(err, "at line:%d", index)
	}
	return nil
}

type RowCallback func(row []string) error

func handleRows(reader *csv.Reader, callback RowCallback) (int64, error) {
	var err error
	var index int64 = 0
	for {
		var row []string
		row, err = reader.Read()
		if err != nil && err != io.EOF {
			logger.Error.Printf("fail to read csv:\n%v", err)
			break
		}
		if err == io.EOF {
			err = nil //EOF是预期的，无需处理
			break
		}
		index++
		//回调处理函数
		err = callback(row)
		if err != nil { //有错误就退出
			break
		}
	}
	logger.Info.Printf("finish handling all rows")
	return index, err
}

// handleRows的回调
func (handler CSVToESHandler) addToES(row []string) error {
	var esMarc = model.ESBookMarc{
		MARCRecNo: row[0],
		CallNo:    row[1],
		Title:     row[2],
		Author:    row[3],
		Publisher: row[4],
	}
	return handler.bulkSession.Add(esMarc.MARCRecNo, esMarc)
}
