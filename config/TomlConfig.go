package config

import (
	"github.com/pelletier/go-toml"
	"library-manager/logger"
)

var sInstanceMap = make(map[string]*TConfig)

type TConfig struct {
	mRootTree *toml.Tree
}

func GetConfig(alias string) *TConfig {
	var config = sInstanceMap[alias]
	if config == nil {
		logger.Error.Fatalln("must call InitConfig first")
		return nil
	}
	return config
}

func InitConfig(alias string, configPath string) {
	rootTree, err := toml.LoadFile(configPath)
	if err != nil {
		logger.Error.Fatalln(err)
		return
	}
	sInstanceMap[alias] = &TConfig{mRootTree: rootTree}
}

func (config TConfig) GetRootTree() *toml.Tree {
	return config.mRootTree
}

func (config TConfig) GetAsTree(key string) *toml.Tree {
	return config.mRootTree.Get(key).(*toml.Tree)
}
