package entities

import (
	"io/ioutil"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/dbtest"
)

func TestInsert(t *testing.T) {
	Convey("Test user insertion", t, func() {
		WrapperUser(func() {
			usr := getAdminUser()
			usr.Insert()
			usrCompare, err := UserFromId(bson.ObjectIdHex("555a076a2fd06c1891000002"))
			Convey("usr to compare should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("User Nickname should be Admin", func() {
				So(usrCompare.Name.NickName, ShouldEqual, usr.Name.NickName)
			})
		})
	})
}

func WrapperUser(toWrap func()) {
	var server dbtest.DBServer
	defer server.Stop()
	tempDir, err := ioutil.TempDir("", "testing")
	Convey("TempDir err should be nil", func() {
		So(err, ShouldBeNil)
	})
	server.SetPath(tempDir)
	session := server.Session()
	db := session.DB("test_mewpipe")
	defer session.Close()
	Convey("db should be nil", func() {
		So(err, ShouldBeNil)
	})
	ChangeUserDB(db)
	toWrap()
}

func getAdminUser() *User {
	usr := UserNewFromId(bson.ObjectIdHex("555a076a2fd06c1891000002"))
	usr.Email = "admin@admin.com"
	usr.Password = "Admin"
	usr.Name.FirstName = "Admin"
	usr.Name.LastName = "Admin"
	usr.Name.NickName = "Admin"
	usr.Roles = append(usr.Roles, "Admin")
	return usr
}
