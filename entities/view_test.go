package entities

import (
	"io/ioutil"
	"supinfo/mewpipe/configs"

	. "github.com/smartystreets/goconvey/convey"

	"testing"
)

func init() {

	tempDir, _ := ioutil.TempDir("", "testing")
	Server.SetPath(tempDir)

	Session = Server.Session()
	configs.MongoDB = Session.DB("test_mewpipe")
}

func TestViewNewAnonymous(t *testing.T) {
	Convey("Test new Anonymous view", t, func() {
		Wipe()
		media := getAmazingMedia()
		media.Insert()
		ViewNewAnonymous(media.Id)
		//Can't be writing unless get Views from a media
	})
}

func TestViewNew(t *testing.T) {
	Convey("Test new view", t, func() {
		Wipe()
		media := getAmazingMedia()
		user := getFooUser()
		user.Insert()
		media.Insert()
		ViewNew(user.Id, media.Id)
		//Can't be writing unless get Views from a media
	})
}
