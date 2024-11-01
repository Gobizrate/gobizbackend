package config

import (
	"github.com/gocroot/helper/atdb"
)

var MongoString string = "mongodb+srv://abnormal:11akusayangibu@gobizcroot.o8vp5.mongodb.net/"

var mongoinfo = atdb.DBInfo{
	DBString: MongoString,
	DBName:   "gobizdevlop",
}

var Mongoconn, ErrorMongoconn = atdb.MongoConnect(mongoinfo)
