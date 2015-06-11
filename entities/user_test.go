package entities

import (
	"io/ioutil"
	"supinfo/mewpipe/configs"
	"testing"

	"encoding/base64"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

func TestUserInsert(t *testing.T) {
	Wipe()
	Convey("Test user insertion", t, func() {
		usr := getFooUser()
		usr.Insert()
		usrCompare, err := UserFromId(bson.ObjectIdHex("5578b8c4f711886e75dec3fd"))
		Convey("User should be found", func() {
			So(err, ShouldBeNil)
		})
		Convey("User Nickname should be Foo", func() {
			So(usrCompare.Name.NickName, ShouldEqual, usr.Name.NickName)
		})
	})
}

func TestUserDelete(t *testing.T) {
	Wipe()
	Convey("Test user removal", t, func() {
		usr := getFooUser()
		usr.Insert()
		_, err := UserFromId(bson.ObjectIdHex("5578b8c4f711886e75dec3fd"))
		Convey("User should be found", func() {
			So(err, ShouldBeNil)
		})
		usr.Delete()
		_, err = UserFromId(bson.ObjectIdHex("5578b8c4f711886e75dec3fd"))
		Convey("User should be not found", func() {
			So(err, ShouldNotBeNil)
		})
	})
}

func TestUserTokenGeneration(t *testing.T) {
	Wipe()
	Convey("Test user token generation", t, func() {
		usr := getFooUser()
		usr.Insert()
		_, err := UserFromId(bson.ObjectIdHex("5578b8c4f711886e75dec3fd"))
		Convey("User should be found", func() {
			So(err, ShouldBeNil)
		})
		token, err := usr.TokenNew()
		_, err = UserFromToken(base64.StdEncoding.EncodeToString([]byte(token.Token)))
		Convey("User should be found from token", func() {
			So(err, ShouldBeNil)
		})
	})
}

func getFooUser() *User {
	usr := UserNewFromId(bson.ObjectIdHex("5578b8c4f711886e75dec3fd"))
	usr.Email = "foo@bar.tld"
	usr.Password = "Foo"
	usr.Name.FirstName = "Foo"
	usr.Name.LastName = "Foo"
	usr.Name.NickName = "Foo"
	return usr
}
