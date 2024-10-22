package config

import (
	"github.com/gocroot/helper/atdb"
)

var MongoString string = "mongodb+srv://ayalarifki:hHyX4lN7TmBtXW38@cluster0.5p1ozyb.mongodb.net/"

var mongoinfo = atdb.DBInfo{
	DBString: MongoString,
	DBName:   "gobizdev",
}

var Mongoconn, ErrorMongoconn = atdb.MongoConnect(mongoinfo)
