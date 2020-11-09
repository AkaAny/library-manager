package orm

import (
	"github.com/pelletier/go-toml"
	"gorm.io/gorm"
	"library-manager/config"
	"library-manager/logger"
)

type BotGORM struct {
	mConfig DBConfig
	mDB     *gorm.DB
}

func init() {
	config.InitConfig(DB_CONFIG_ALIAS, DB_CONFIG_PATH)
}

func GetDBConfig() *config.TConfig {
	return config.GetConfig(DB_CONFIG_ALIAS)
}

func MustCreateLibGOrm(moduleName string) *BotGORM {
	var dbConfig DBConfig
	moduleTree := GetDBConfig().GetRootTree().Get(moduleName).(*toml.Tree)
	dbTree := moduleTree.Get("database").(*toml.Tree)
	err := dbTree.Unmarshal(&dbConfig)
	if err != nil {
		logger.Error.Fatalln(err)
		return nil
	}
	var dbType = dbConfig.Type
	var creator = GetCreatorByType(dbType)
	if creator == nil {
		logger.Error.Fatalf("fail to find creator for type:%s", dbType)
		return nil
	}
	db, err := creator.Create(dbConfig)
	if err != nil {
		logger.Error.Fatalln(err)
		return nil
	}

	return &BotGORM{mConfig: dbConfig, mDB: db}
}

func (sgo BotGORM) GetDB() *gorm.DB {
	return sgo.mDB
}

func (sgo BotGORM) IsModelExist(model interface{}, out interface{}) bool {
	var count int64
	sgo.GetDB().First(out).Count(&count)
	return count > 0
}

func (sgo BotGORM) PutModel(model interface{}, out interface{}) {
	var exist = sgo.IsModelExist(model, out)
	if exist {
		sgo.GetDB().Create(model)
	} else {
		sgo.GetDB().Save(model)
	}
}
