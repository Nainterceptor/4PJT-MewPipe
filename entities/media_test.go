package entities

import (
	"io/ioutil"
	"os"

	"github.com/Nainterceptor/4PJT-MewPipe/configs"

	"time"

	"testing"

	"mime/multipart"

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
			//Upload Media
			header := new(multipart.FileHeader)
			header.Filename = "sample.mp4"
			header.Header = make(map[string][]string)
			header.Header.Add("Content-Type", "video/mp4")
			sampleFile, err := os.Open("../fixtures/files/sample.mp4")
			Convey("Sample reading must be a success", func() {
				So(err, ShouldBeNil)
			})
			err = media.Upload(sampleFile, header)
			Convey("Sample upload should be a success", func() {
				So(err, ShouldBeNil)
			})

			//Upload Thumbnail
			headerThumb := new(multipart.FileHeader)
			headerThumb.Filename = "sample.jpg"
			headerThumb.Header = make(map[string][]string)
			headerThumb.Header.Add("Content-Type", "image/jpeg")
			sampleThumb, err := os.Open("../fixtures/files/sample.jpg")
			Convey("Thumbnail reading must be a success", func() {
				So(err, ShouldBeNil)
			})
			err = media.UploadThumbnail(sampleThumb, headerThumb)
			Convey("Thumbnail upload should be a success", func() {
				So(err, ShouldBeNil)
			})
			Convey("All files should been deleted", func() {
				So(countFiles(), ShouldEqual, 1)
				So(countChunks(), ShouldAlmostEqual, 1)
				So(countThumbnailFiles(), ShouldEqual, 1)
				So(countThumbnailChunks(), ShouldEqual, 1)
				So(media.Delete(), ShouldBeNil)
				So(countFiles(), ShouldEqual, 0)
				So(countChunks(), ShouldEqual, 0)
				So(countThumbnailFiles(), ShouldEqual, 0)
				So(countThumbnailChunks(), ShouldEqual, 0)
			})
		})

	})
}

func TestMediaUpload(t *testing.T) {
	Convey("Test media upload", t, func() {
		Wipe()
		media := getAmazingMedia()
		media.Insert()
		header := new(multipart.FileHeader)
		header.Filename = "sample.mp4"
		header.Header = make(map[string][]string)
		header.Header.Add("Content-Type", "video/mp4")
		sampleFile, err := os.Open("../fixtures/files/sample.mp4")
		Convey("Sample reading must be a success", func() {
			So(err, ShouldBeNil)
		})
		err = media.Upload(sampleFile, header)
		Convey("Sample upload should be a success", func() {
			So(err, ShouldBeNil)
		})
		Convey("Sample upload should be a success again", func() {
			So(countFiles(), ShouldEqual, 1)
			err = media.Upload(sampleFile, header)
			So(err, ShouldBeNil)
			So(countFiles(), ShouldEqual, 1)
		})
		Convey("Sample upload should fail if Update fail", func() {
			media.Id = bson.NewObjectId()
			err = media.Upload(sampleFile, header)
			So(err, ShouldNotBeNil)
		})
		Convey("Sample upload should fail when file is closed", func() {
			sampleFile.Close()
			err = media.Upload(sampleFile, header)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestThumbnailUpload(t *testing.T) {
	Convey("Test thumbnail upload", t, func() {
		Wipe()
		media := getAmazingMedia()
		media.Insert()
		header := new(multipart.FileHeader)
		header.Filename = "sample.jpg"
		header.Header = make(map[string][]string)
		header.Header.Add("Content-Type", "image/jpeg")
		sampleFile, err := os.Open("../fixtures/files/sample.jpg")
		Convey("Sample reading must be a success", func() {
			So(err, ShouldBeNil)
		})
		err = media.UploadThumbnail(sampleFile, header)
		Convey("Sample upload should be a success", func() {
			So(err, ShouldBeNil)
			So(media.Thumbnail, ShouldNotBeEmpty)
		})
		Convey("Sample upload should be a success again", func() {
			So(countThumbnailFiles(), ShouldEqual, 1)
			err = media.UploadThumbnail(sampleFile, header)
			So(err, ShouldBeNil)
			So(countThumbnailFiles(), ShouldEqual, 1)
		})
		Convey("Sample upload should fail if Update fail", func() {
			media.Id = bson.NewObjectId()
			err = media.UploadThumbnail(sampleFile, header)
			So(err, ShouldNotBeNil)
		})
		Convey("Sample upload should fail when file is closed", func() {
			sampleFile.Close()
			err = media.UploadThumbnail(sampleFile, header)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestMediaRead(t *testing.T) {
	Convey("Test media Read", t, func() {
		Wipe()
		media := getAmazingMedia()
		media.Insert()
		header := new(multipart.FileHeader)
		header.Filename = "sample.mp4"
		header.Header = make(map[string][]string)
		header.Header.Add("Content-Type", "video/mp4")
		sampleFile, err := os.Open("../fixtures/files/sample.mp4")
		Convey("Sample reading must be a success", func() {
			So(err, ShouldBeNil)
		})
		err = media.Upload(sampleFile, header)
		Convey("Sample upload should be a success", func() {
			So(err, ShouldBeNil)
		})

		//Now, Read !
		err = media.OpenFile()
		defer media.CloseFile()
		Convey("Media open should be a success", func() {
			So(err, ShouldBeNil)
		})

		Convey("Media should be mp4", func() {
			So(media.ContentType(), ShouldEqual, "video/mp4")
		})

		start := 0
		intSize := int(media.Size())
		err = media.SeekSet(int64(start))
		Convey("Seek should be a success", func() {
			So(err, ShouldBeNil)
		})
		buffer := make([]byte, intSize)

		err = media.Read(buffer)
		Convey("Read should be a success", func() {
			So(err, ShouldBeNil)
		})
	})
}

func TestThumbnailRead(t *testing.T) {
	Convey("Test media Read", t, func() {
		Wipe()
		media := getAmazingMedia()
		media.Insert()
		header := new(multipart.FileHeader)
		header.Filename = "sample.jpg"
		header.Header = make(map[string][]string)
		header.Header.Add("Content-Type", "image/jpeg")
		sampleFile, err := os.Open("../fixtures/files/sample.jpg")
		Convey("Sample reading must be a success", func() {
			So(err, ShouldBeNil)
		})
		err = media.UploadThumbnail(sampleFile, header)
		Convey("Sample upload should be a success", func() {
			So(err, ShouldBeNil)
		})

		//Now, Read !
		err = media.OpenThumbnail()
		defer media.CloseThumbnail()
		Convey("Thumbnail open should be a success", func() {
			So(err, ShouldBeNil)
		})

		intSize := int(media.ThumbnailSize())
		_ = make([]byte, intSize)

		Convey("Thumbnail Read should be a success", func() {
			file, err := ioutil.TempFile("", "fixtures_")
			Convey("Thumbnail Read should be a success", func() {
				So(err, ShouldBeNil)
			})
			So(media.CopyThumbnailTo(file), ShouldBeNil)
		})
	})
}

func TestMediaCopy(t *testing.T) {
	Convey("Test media copy", t, func() {
		Wipe()
		media := getAmazingMedia()
		media.Insert()
		header := new(multipart.FileHeader)
		header.Filename = "sample.mp4"
		header.Header = make(map[string][]string)
		header.Header.Add("Content-Type", "video/mp4")
		sampleFile, err := os.Open("../fixtures/files/sample.mp4")
		Convey("Sample reading must be a success", func() {
			So(err, ShouldBeNil)
		})
		err = media.Upload(sampleFile, header)
		Convey("Sample upload should be a success", func() {
			So(err, ShouldBeNil)
		})

		//Now, Read !
		err = media.OpenFile()
		defer media.CloseFile()
		Convey("Media open should be a success", func() {
			So(err, ShouldBeNil)
		})
		Convey("Media Read should be a success", func() {
			file, err := ioutil.TempFile("", "fixtures_")
			Convey("Media Read should be a success", func() {
				So(err, ShouldBeNil)
			})
			So(media.CopyTo(file), ShouldBeNil)
		})
	})
}

func TestMediaOpen(t *testing.T) {
	Convey("Test media open", t, func() {
		Wipe()
		media := getAmazingMedia()
		Convey("Media without file should back an error", func() {
			So(media.OpenFile(), ShouldNotBeNil)
		})
	})

}

func TestCountViews(t *testing.T) {
	Convey("Test count views", t, func() {
		Wipe()
		media := getAmazingMedia()
		user := getFooUser()
		user2 := getBarUser()
		media.Insert()
		user.Insert()
		user2.Insert()
		// 4 view from anonymous
		for i := 0; i < 4; i++ {
			ViewNewAnonymous(media.Id)
		}
		// 8 view from Foo
		for i := 0; i < 8; i++ {
			ViewNew(user.Id, media.Id)
		}
		// 16 view from Bar
		for i := 0; i < 16; i++ {
			ViewNew(user2.Id, media.Id)
		}

		media.CountViews()
		So(media.Views, ShouldEqual, 28)
	})

}

func TestCountShares(t *testing.T) {
	Convey("Test count shares", t, func() {
		Wipe()
		media := getAmazingMedia()
		user := getFooUser()
		user2 := getBarUser()
		media.Insert()
		user.Insert()
		user2.Insert()
		// 4 view from anonymous
		for i := 0; i < 4; i++ {
			ShareCountNewAnonymous(media.Id)
		}
		// 8 view from Foo
		for i := 0; i < 8; i++ {
			ShareCountNew(user.Id, media.Id)
		}
		// 16 view from Bar
		for i := 0; i < 16; i++ {
			ShareCountNew(user2.Id, media.Id)
		}

		media.CountShares()
		So(media.Shares, ShouldEqual, 28)
	})

}

func countFiles() int {
	count, _ := configs.MongoDB.C("media.files").Count()
	return count
}

func countChunks() int {
	count, _ := configs.MongoDB.C("media.chunks").Count()
	return count
}

func countThumbnailFiles() int {
	count, _ := configs.MongoDB.C("media.thumbnails.files").Count()
	return count
}

func countThumbnailChunks() int {
	count, _ := configs.MongoDB.C("media.thumbnails.chunks").Count()
	return count
}

func getAmazingMedia() *Media {
	media := MediaNew()
	media.Title = "Amazing video"
	media.Summary = "In this video, a big comet hit Earth, and Superman dances with a unicorn."
	user := getFooUser()
	media.Publisher.Name = user.Name
	media.Publisher.Email = user.Email
	media.Publisher.Id = user.Id
	media.Scope = Public

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
	media.Scope = Private

	return media
}
