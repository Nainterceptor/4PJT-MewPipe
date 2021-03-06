package entities

import (
	"io/ioutil"
	"testing"

	"github.com/Nainterceptor/4PJT-MewPipe/configs"

	"encoding/base64"

	"time"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2/bson"
)

func init() {

	tempDir, _ := ioutil.TempDir("", "testing")
	Server.SetPath(tempDir)

	Session = Server.Session()
	configs.MongoDB = Session.DB("test_mewpipe")
}

func TestUserNew(t *testing.T) {
	usr := UserNew()
	Convey("Test UserNew has informations", t, func() {
		So(usr.Id, ShouldNotBeEmpty)
	})
}

func TestUserTokenNew(t *testing.T) {
	token := userTokenNew()
	Convey("Test userToken has informations", t, func() {
		So(token.Token, ShouldNotBeEmpty)
		currentTime := time.Now()
		So(token.ExpireAt.Second(), ShouldEqual, currentTime.Add(time.Second*tokenExpiration).Second()) //Precision on second
	})
}

func TestUserNewFromId(t *testing.T) {
	id := bson.ObjectIdHex("5578b8c4f711886e75dec3fd")
	usr := UserNewFromId(id)
	Convey("Test UserNew has right ID", t, func() {
		So(usr.Id, ShouldEqual, bson.ObjectIdHex("5578b8c4f711886e75dec3fd"))
	})
}

func TestUserList(t *testing.T) {
	Wipe()
	Convey("Test list users", t, func() {
		getFooUser().Insert()
		getBarUser().Insert()
		getAdminUser().Insert()
		users, err := UserList(bson.M{}, 0, 10)
		Convey("UserList should not have an error", func() {
			So(err, ShouldBeNil)
		})
		Convey("UserList should have 3 results", func() {
			So(len(users), ShouldEqual, 3)
		})
		users2, err := UserList(bson.M{}, 0, 2)
		Convey("UserList 2 should not have an error", func() {
			So(err, ShouldBeNil)
		})
		Convey("UserList should have 2 results", func() {
			So(len(users2), ShouldEqual, 2)
		})
		users3, err := UserList(bson.M{}, 1, 1)
		Convey("UserList 3 should not have an error", func() {
			So(err, ShouldBeNil)
		})
		Convey("UserList should be same", func() {
			So(users3[:1][0].Id, ShouldEqual, users[1:2][0].Id)
		})
	})
}

func TestUserValidation(t *testing.T) {
	Convey("Test user validation", t, func() {
		Convey("Fullfilled user should not back an error", func() {
			usr := getFooUser()
			So(usr.Validate(), ShouldBeNil)
		})
		Convey("User email should not be empty if twitter is empty", func() {
			usr := getFooUser()
			usr.Email = ""
			So(usr.Validate(), ShouldNotBeNil)
		})
		Convey("User email should be empty if twitter exist", func() {
			usr := getFooUser()
			usr.Email = ""
			usr.Twitter.UserId = "99999"
			So(usr.Validate(), ShouldBeNil)
		})
		Convey("User email should be valid if twitter is empty", func() {
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
		media := getAmazingMedia()
		media.Insert()
		usr.Delete()
		_, err = UserFromId(bson.ObjectIdHex("5578b8c4f711886e75dec3fd"))
		Convey("User should not be found", func() {
			So(err, ShouldNotBeNil)
		})
		usr = getBarUser()
		Convey("Removal on inexisting ID should fail", func() {
			So(usr.Delete(), ShouldNotBeNil)
		})
	})
}

func TestUserGetMedia(t *testing.T) {
	usr := getFooUser()
	usr.GetMedia()
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
		getAmazingMedia().Insert()
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
		Convey("Media must be updated cascade", func() {
			media, _ := usr.GetMedia()
			So(media[0].Publisher.Email, ShouldEqual, usr.Email)
		})
		usr = getBarUser()
		Convey("Update on inexisting ID should fail", func() {
			So(usr.Update(), ShouldNotBeNil)
		})

	})
}

func TestUserFromCredentials(t *testing.T) {
	Wipe()
	Convey("Test get user from credentials", t, func() {
		usr := getFooUser()
		usr.Insert()
		_, err := UserFromCredentials("foo@bar.tld", "Foo")
		Convey("User should be found", func() {
			So(err, ShouldBeNil)
		})
		_, err = UserFromCredentials("foo@bar.tld", "BadPassword")
		Convey("User should not be found", func() {
			So(err, ShouldNotBeNil)
		})
	})
}

func TestUserFromTwitter(t *testing.T) {
	Wipe()
	Convey("Test get user from twitter", t, func() {
		usr := getTwitterUser()
		usr.Insert()
		_, err := UserFromTwitterUserID("Twitter")
		Convey("User should be found", func() {
			So(err, ShouldBeNil)
		})
		_, err = UserFromTwitterUserID("BadTwitter")
		Convey("User should not be found", func() {
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
		tokenBase64 := base64.StdEncoding.EncodeToString([]byte(token.Token))
		_, err = UserFromToken(tokenBase64)
		Convey("User must have token", func() {
			So(usr.hasToken(string(token.Token)), ShouldBeTrue)
		})
		Convey("User should be found from token", func() {
			So(err, ShouldBeNil)
		})
		usr = getBarUser()
		_, err = usr.TokenNew()
		Convey("Update on inexisting ID should fail", func() {
			So(err, ShouldNotBeNil)
		})
	})
}

func TestUserFromToken(t *testing.T) {
	Wipe()
	Convey("Test user from token", t, func() {
		usr := getFooUser()
		usr.Insert()
		token, _ := usr.TokenNew()
		tokenBase64 := base64.StdEncoding.EncodeToString([]byte(token.Token))
		_, err := UserFromToken("BadToken$@`£")
		Convey("Incorrect token should return an error", func() {
			So(err, ShouldNotBeNil)
		})
		_, err = UserFromToken("BadToken")
		Convey("Bad token should return an error", func() {
			So(err, ShouldNotBeNil)
		})
		_, err = UserFromToken(tokenBase64)
		Convey("Good token should not return an error", func() {
			So(err, ShouldBeNil)
		})
		usr.UserTokens[0].ExpireAt = time.Now().Add(-1*time.Second*tokenExpiration - 1)
		//Do not use Update, which clean the old token.
		getUserCollection().UpdateId(usr.Id, &usr)
		_, err = UserFromToken(tokenBase64)
		Convey("Old token should return an error", func() {
			So(err, ShouldNotBeNil)
		})
	})
}

func TestUserHasRole(t *testing.T) {
	Wipe()
	Convey("Test user has role", t, func() {
		usr := getFooUser()
		Convey("User should not have Test role", func() {
			So(usr.HasRole("Test"), ShouldBeFalse)
		})
		usr.Roles = append(usr.Roles, "Test")
		Convey("User should have Test role", func() {
			So(usr.HasRole("Test"), ShouldBeTrue)
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

func getBarUser() *User {
	usr := UserNewFromId(bson.ObjectIdHex("557975fb2ca00357367f7a98"))
	usr.Email = "bar@foo.tld"
	usr.Password = "Bar"
	usr.Name.FirstName = "Bar"
	usr.Name.LastName = "Bar"
	usr.Name.NickName = "Bar"
	return usr
}

func getAdminUser() *User {
	usr := UserNewFromId(bson.ObjectIdHex("557976072ca00357367f7a99"))
	usr.Email = "Admin@admin.tld"
	usr.Password = "Admin"
	usr.Name.FirstName = "Admin"
	usr.Name.LastName = "Admin"
	usr.Name.NickName = "Admin"
	usr.Roles = append(usr.Roles, "Admin")
	return usr
}

func getTwitterUser() *User {
	usr := UserNew()
	usr.Name.NickName = "Twitter"
	usr.Twitter.UserId = "Twitter"
	return usr
}
