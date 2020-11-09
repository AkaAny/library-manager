package es

import (
	"github.com/elastic/go-elasticsearch/v8"
	"library-manager/config"
	"library-manager/logger"
)

type LibraryES struct {
	mClient *elasticsearch.Client
}

func (les LibraryES) GetClient() *elasticsearch.Client {
	return les.mClient
}

func MustCreateFromConfig() *LibraryES {
	config.InitConfig(CONFIG_ALIAS, CONFIG_PATH)
	var esTree = config.GetConfig(CONFIG_ALIAS).GetAsTree("elasticsearch")
	var esConfig ESConfig
	err := esTree.Unmarshal(&esConfig)
	if err != nil {
		logger.Error.Fatalln(err)
		return nil
	}
	logger.Info.Printf("es address:%v", esConfig.Addresses)
	client, err := elasticsearch.NewClient(esConfig.ToClientConfig())
	if err != nil {
		logger.Error.Fatalln(err)
		return nil
	}
	return &LibraryES{mClient: client}
}
