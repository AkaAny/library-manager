package dbimport

import (
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"library-manager/model"
	"strings"
)

type BulkSession struct {
	Indexer esutil.BulkIndexer
}

func (bs BulkSession) GetIndexer() esutil.BulkIndexer {
	return bs.Indexer
}

type ImportError struct {
	Item     esutil.BulkIndexerItem
	Response esutil.BulkIndexerResponseItem
	Error    error
}

func AppendWithError(fails []ImportError, err error) []ImportError {
	fails = append(fails, ImportError{Error: err})
	return fails //切片作为参数时是值传递
}

func (bs *BulkSession) Add(docId string, bodyObj fmt.Stringer) error {
	var bulkItem = esutil.BulkIndexerItem{
		Action:     "index",
		DocumentID: docId,
		Body:       strings.NewReader(bodyObj.String()),
	}
	return bs.GetIndexer().Add(context.Background(), bulkItem)
}

func (bs BulkSession) Flush() error {
	return bs.GetIndexer().Close(context.Background())
}

func Create(client *elasticsearch.Client, index string) (*BulkSession, error) {
	indexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client:  client,
		Index:   index,
		Refresh: "true",
	})
	if err != nil {
		return nil, err
	}
	return &BulkSession{Indexer: indexer}, nil
}

func doImport(client *elasticsearch.Client, marcs []model.ESBookMarc) []ImportError {
	var fails []ImportError
	indexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client:  client,
		Index:   "marc",
		Refresh: "true",
	})
	if err != nil {
		fails = AppendWithError(fails, err)
		return fails
	}
	for _, marc := range marcs {
		var bulkItem = esutil.BulkIndexerItem{
			Action:     "index",
			DocumentID: marc.MARCRecNo,
			Body:       strings.NewReader(marc.ToJSON()),
			OnFailure: func(ctx context.Context,
				item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem,
				err error) {
				fails = append(fails, ImportError{
					Item:     item,
					Response: res,
					Error:    err,
				})
			},
		}
		err = indexer.Add(context.Background(), bulkItem)
		if err != nil {
			fails = append(fails, ImportError{
				Error: err,
			})
			break
		}
	}
	err = indexer.Close(context.Background())
	if err != nil {
		fails = AppendWithError(fails, err)
		return fails
	}
	return fails
}
