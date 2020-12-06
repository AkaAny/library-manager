package es

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"library-manager/config"
	"library-manager/logger"
)

type LibraryES struct {
	mClient *elasticsearch.Client
}

func (les LibraryES) GetClient() *elasticsearch.Client {
	return les.mClient
}

func (les LibraryES) Search(index string, cond ESBool) (*esapi.Response, error) {
	var body = H{
		"query": H{
			"bool": cond,
		},
	}
	rawData, err := json.Marshal(body)
	logger.Info.Printf("search body:\n%s", rawData)

	if err != nil {
		return nil, err
	}
	var client = les.GetClient()
	var search = client.Search
	return client.Search(
		search.WithContext(context.Background()),
		search.WithIndex(index),
		search.WithBody(bytes.NewReader(rawData)),
		search.WithTrackTotalHits(true),
		search.WithPretty(),
	)
}

func (les LibraryES) SearchBySQL(sql string) {
	query := map[string]interface{}{
		"query":      sql,
		"fetch_size": 5,
	}
	jsonBody, _ := json.Marshal(query)
	var queryRequest = esapi.SQLQueryRequest{
		Body:   bytes.NewReader(jsonBody),
		Pretty: true,
	}
	resp, err := queryRequest.Do(context.Background(), les.GetClient())
	if err != nil {
		logger.Error.Fatalln(err)
	}
	logger.Info.Println(resp.String())
}

func CreateFromConfig() (*LibraryES, error) {
	config.InitConfig(CONFIG_ALIAS, CONFIG_PATH)
	var esTree = config.GetConfig(CONFIG_ALIAS).GetAsTree("elasticsearch")
	var esConfig ESConfig
	err := esTree.Unmarshal(&esConfig)
	if err != nil {
		return nil, err
	}
	logger.Info.Printf("es address:%v", esConfig.Addresses)
	client, err := elasticsearch.NewClient(esConfig.ToClientConfig())
	if err != nil {
		return nil, err
	}
	return &LibraryES{mClient: client}, nil
}

func (les LibraryES) GetInfo() (*esapi.Response, error) {
	return les.GetClient().Info()
}
