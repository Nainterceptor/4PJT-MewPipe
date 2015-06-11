package entities

import (
	"io/ioutil"
	"supinfo/mewpipe/configs"
	"testing"

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

func TestMediaNew(t *testing.T) {
	Convey("Test Media new", t, func() {
		media := MediaNew()
		Convey("Test default var", func() {
			So(media.Id, ShouldNotBeEmpty)
			So(media.CreatedAt.Second(), ShouldEqual, time.Now().Second())
		})
	})
}

func TestMediaNewFromId(t *testing.T) {
	Convey("Test Media new from ID", t, func() {
		oid := bson.ObjectIdHex("5579d805cd46920dbcaf2ff8")
		media := MediaNewFromId(oid)
		Convey("Test media ID", func() {
			So(media.Id, ShouldEqual, oid)
			So(media.CreatedAt.Second(), ShouldEqual, time.Now().Second())
		})
	})
}

func TestMediaFromId(t *testing.T) {
	Convey("Test Media from ID", t, func() {
		media := getAmazingMedia()
		media.Id = bson.ObjectIdHex("5579d805cd46920dbcaf2ff8")
		_, err := MediaFromId(media.Id)
		Convey("Media not exist should return an error", func() {
			So(err, ShouldNotBeNil)
		})
		media.Insert()
		_, err = MediaFromId(media.Id)
		Convey("Media which exist should not an error", func() {
			So(err, ShouldBeNil)
		})
	})
}

func TestMediaList(t *testing.T) {
	Convey("Test list media", t, func() {
		Wipe()
		getAmazingMedia().Insert()
		getBadMedia().Insert()
		media, err := MediaList(bson.M{}, 0, 10)
		Convey("MediaList should not have an error", func() {
			So(err, ShouldBeNil)
		})
		Convey("MediaList should have 2 results", func() {
			So(len(media), ShouldEqual, 2)
		})
		media2, err := MediaList(bson.M{}, 0, 1)
		Convey("MediaList 2 should not have an error", func() {
			So(err, ShouldBeNil)
		})
		Convey("MediaList should have 1 results", func() {
			So(len(media2), ShouldEqual, 1)
		})
		media3, err := MediaList(bson.M{}, 1, 1)
		Convey("MediaList 3 should not have an error", func() {
			So(err, ShouldBeNil)
		})
		Convey("MediaList should be same", func() {
			So(media3[:1][0].Id, ShouldEqual, media[1:2][0].Id)
		})
	})
}

func TestMediaNormalize(t *testing.T) {
	Convey("Test media normalize", t, func() {
		media := getAmazingMedia()
		Convey("Test normalize scope", func() {
			media.Scope = "link"
			media.Normalize()
			So(media.Scope, ShouldEqual, "link")
			media.Scope = "foo"
			media.Normalize()
			So(media.Scope, ShouldEqual, "public")
		})
	})
}

func TestMediaInsert(t *testing.T) {
	Convey("Test Media insertion", t, func() {
		Wipe()
		media := getAmazingMedia()
		So(media.Insert(), ShouldBeNil)
		So(media.Insert(), ShouldNotBeNil)
	})
}

func TestMediaUpdate(t *testing.T) {
	Convey("Test Media updating", t, func() {
		Wipe()
		media := getAmazingMedia()
		Convey("Test update not inserted", func() {
			So(media.Update(), ShouldNotBeNil)
		})
		Convey("Test update", func() {
			media.Insert()
			So(media.Update(), ShouldBeNil)
		})
	})
}

func TestMediaDelete(t *testing.T) {
	Convey("Test Media removal", t, func() {
		Wipe()
		media := getAmazingMedia()
		Convey("Test delete not inserted", func() {
			So(media.Delete(), ShouldNotBeNil)
		})
		Convey("Test delete", func() {
			media.Insert()
			So(media.Delete(), ShouldBeNil)
		})
	})
}

func getAmazingMedia() *Media {
	media := MediaNew()
	media.Title = "Amazing video"
	media.Summary = "In this video, a big comet hit Earth, and Superman dances with a unicorn."
	user := getFooUser()
	media.Publisher.Name = user.Name
	media.Publisher.Email = user.Email
	media.Publisher.Id = user.Id
	media.Scope = "public"

	return media
}

func getBadMedia() *Media {
	media := MediaNew()
	media.Title = "Another lolcat video"
	media.Summary = "In this video, a lolcat play with catsnip !"
	user := getFooUser()
	media.Publisher.Name = user.Name
	media.Publisher.Email = user.Email
	media.Publisher.Id = user.Id
	media.Scope = "private"

	return media
}
