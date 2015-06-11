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

func TestUserValidation(t *testing.T) {
	Convey("Test user validation", t, func() {
		Convey("Fullfilled user should not back an error", func() {
			usr := getFooUser()
			So(usr.Validate(), ShouldBeNil)
		})
		Convey("User email should not be empty", func() {
			usr := getFooUser()
			usr.Email = ""
			So(usr.Validate(), ShouldNotBeNil)
		})
		Convey("User email should be valid", func() {
			usr := getFooUser()
			usr.Email = "NotAnEmail"
			So(usr.Validate(), ShouldNotBeNil)
		})
		Convey("User Nickname should not be empty", func() {
			usr := getFooUser()
			usr.Name.NickName = ""
			So(usr.Validate(), ShouldNotBeNil)
		})
	})
}

func TestUserNormalize(t *testing.T) {
	Convey("Test user normalize", t, func() {
		Convey("User must be trimed", func() {
			usr := getFooUser()
			usr.Email = " " + usr.Email + " "
			usr.Name.NickName = " " + usr.Name.NickName + " "
			usr.Name.LastName = " " + usr.Name.LastName + " "
			usr.Name.FirstName = " " + usr.Name.FirstName + " "
			usr.Normalize()
			So(usr.Email, ShouldEqual, "foo@bar.tld")
			So(usr.Name.NickName, ShouldEqual, "Foo")
			So(usr.Name.LastName, ShouldEqual, "Foo")
			So(usr.Name.FirstName, ShouldEqual, "Foo")
		})

	})
}

func TestUserInsert(t *testing.T) {
	Wipe()
	Convey("Test user insertion", t, func() {
		usr := getFooUser()
		usr.Insert()
		_, err := UserFromId(bson.ObjectIdHex("5578b8c4f711886e75dec3fd"))
		Convey("User should be found", func() {
			So(err, ShouldBeNil)
		})
		Convey("User Password should be Empty", func() {
			So(usr.Password, ShouldBeEmpty)
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
		Convey("User should not be found", func() {
			So(err, ShouldNotBeNil)
		})
	})
}

func TestUserUpdate(t *testing.T) {
	Wipe()
	Convey("Test user updating", t, func() {
		usr := getFooUser()
		usr.Insert()
		_, err := UserFromCredentials("foo@bar.tld", "Foo")
		Convey("User should be found", func() {
			So(err, ShouldBeNil)
		})
		usr.Email = "bar@foo.tld"
		usr.Password = "Bar"
		usr.Update()
		_, err = UserFromCredentials("foo@bar.tld", "Foo")
		Convey("Old user should not be found", func() {
			So(err, ShouldNotBeNil)
		})
		_, err = UserFromCredentials("bar@foo.tld", "Bar")
		Convey("New user should be found", func() {
			So(err, ShouldBeNil)
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
