package entities

import (
	"io/ioutil"

	"github.com/Nainterceptor/4PJT-MewPipe/configs"

	. "github.com/smartystreets/goconvey/convey"

	"testing"
)

func init() {

	tempDir, _ := ioutil.TempDir("", "testing")
	Server.SetPath(tempDir)

	Session = Server.Session()
	configs.MongoDB = Session.DB("test_mewpipe")
}

func TestShareCountNewAnonymous(t *testing.T) {
	Convey("Test new Anonymous share", t, func() {
		Wipe()
		media := getAmazingMedia()
		media.Insert()
		ShareCountNewAnonymous(media.Id)
		//Can't be writing unless get Views from a media
	})
}

func TestShareCountNew(t *testing.T) {
	Convey("Test new share", t, func() {
		Wipe()
		media := getAmazingMedia()
		user := getFooUser()
		user.Insert()
		media.Insert()
		ShareCountNew(user.Id, media.Id)
	})
}
