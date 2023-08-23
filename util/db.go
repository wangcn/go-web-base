package util

import (
	"go.uber.org/zap"
	"gorm.io/gorm"

	"mybase/util/database"
)

var dbManager = database.CreateClientManager()

type DBName string

const (
	DBMaster DBName = "db.master"
	DBSlave  DBName = "db.slave"
)

func DB(name DBName) *gorm.DB {

	log := Log().With(zap.String("name", string(name)))
	config := new(database.Config)
	dbName := string(name)

	key := dbName
	if dbEngine := dbManager.Get(key); dbEngine != nil {
		return dbEngine.GetORM()
	}

	err := Cfg("app").UnmarshalKey(dbName, config)
	if err != nil {
		log.With(zap.Error(err)).Error("db load def config err")
		return nil
	}

	client, err := database.CreateClient(key, config)
	if err != nil {
		log.With(zap.Error(err)).Error("db create db engine err")
		return nil
	}
	dbManager.Set(key, client)

	return client.GetORM()
}
