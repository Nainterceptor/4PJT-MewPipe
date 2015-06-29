package entities

import (
	"github.com/Nainterceptor/4PJT-MewPipe/configs"

	"io/ioutil"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/dbtest"
)

var Server dbtest.DBServer
var Session *mgo.Session

func init() {

	tempDir, _ := ioutil.TempDir("", "testing")
	Server.SetPath(tempDir)

	Session = Server.Session()
	configs.MongoDB = Session.DB("test_mewpipe")
}

func Wipe() {
	configs.MongoDB.DropDatabase()
}
