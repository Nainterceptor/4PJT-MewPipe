package entities

import (
	"fmt"
	"math/rand"
	"mime/multipart"
	"net/textproto"
	"os"
	"supinfo/mewpipe/configs"
	"supinfo/mewpipe/entities"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func ClearMedia() {
	configs.MongoDB.C("media.chunks").DropCollection()
	configs.MongoDB.C("media.thumbnails.chunks").DropCollection()
	configs.MongoDB.C("media.files").DropCollection()
	configs.MongoDB.C("media.thumbnails.files").DropCollection()
	configs.MongoDB.C("media").DropCollection()
	configs.MongoDB.C("media.views").DropCollection()
	configs.MongoDB.C("media.shareCounts").DropCollection()
	fmt.Println("All media deleted")
}

func InsertSomeMedia() {
	var mediaArray []*entities.Media
	mediaArray = append(mediaArray, getBadMedia())
	mediaArray = append(mediaArray, getAmazingMedia())
	mediaArray = append(mediaArray, getFooPublicMedia())
	mediaArray = append(mediaArray, getFooPrivateMedia())
	mediaArray = append(mediaArray, getFooLinkMedia())
	mediaArray = append(mediaArray, getAdminPublicMedia())
	mediaArray = append(mediaArray, getAdminPrivateMedia())
	mediaArray = append(mediaArray, getAdminLinkMedia())
	for _, media := range mediaArray {
		if err := media.Insert(); err != nil {
			panic(err)
			return
		}
		addVideoBin(media)
		numAnonViews := randInt(10, 100)
		numFooViews := randInt(0, 7)
		numAdminViews := randInt(0, 3)
		for i := 0; i < numAnonViews; i++ {
			entities.ViewNewAnonymous(media.Id)
		}
		for i := 0; i < numFooViews; i++ {
			entities.ViewNew(bson.ObjectIdHex("555a076a2fd06c1891000001"), media.Id)
		}
		for i := 0; i < numAdminViews; i++ {
			entities.ViewNew(bson.ObjectIdHex("555a076a2fd06c1891000002"), media.Id)
		}
		media.CountViews()
		numAnonShares := randInt(0, 10)
		numFooShares := randInt(0, 3)
		numAdminShares := randInt(0, 1)
		for i := 0; i < numAnonShares; i++ {
			entities.ShareCountNewAnonymous(media.Id)
		}
		for i := 0; i < numFooShares; i++ {
			entities.ShareCountNew(bson.ObjectIdHex("555a076a2fd06c1891000001"), media.Id)
		}
		for i := 0; i < numAdminShares; i++ {
			entities.ShareCountNew(bson.ObjectIdHex("555a076a2fd06c1891000002"), media.Id)
		}
		media.CountShares()
	}
	addVideoBin(mediaArray[1])
	fmt.Println("Media added")
	fmt.Println("Warning : Our sample is packaged to be light (~100ko for 00:14 duration), and contain videos without audio.")
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func addVideoBin(media *entities.Media) {
	//Video
	in, err := os.Open("fixtures/files/sample.mp4")
	if err != nil {
		fmt.Println(err)
		return
	}
	fileHeader := new(multipart.FileHeader)
	fileHeader.Filename = "sample.mp4"
	fileHeader.Header = textproto.MIMEHeader{
		"Content-Type": {"video/mp4; charset=UTF-8"},
	}
	if err = media.Upload(in, fileHeader); err != nil {
		fmt.Println(err)
		return
	}

	//Thumbnail
	thumb, err := os.Open("fixtures/files/sample.jpg")
	if err != nil {
		fmt.Println(err)
		return
	}
	thumbHeader := new(multipart.FileHeader)
	thumbHeader.Filename = "sample.jpg"
	thumbHeader.Header = textproto.MIMEHeader{
		"Content-Type": {"video/mp4; charset=UTF-8"},
	}
	if err = media.UploadThumbnail(thumb, thumbHeader); err != nil {
		fmt.Println(err)
		return
	}
}

func getAmazingMedia() *entities.Media {
	media := entities.MediaNew()
	media.Title = "Amazing video"
	media.Summary = "In this video, a big comet hit Earth, and Superman dances with a unicorn."
	user := getFooUser()
	media.Publisher.Name = user.Name
	media.Publisher.Email = user.Email
	media.Publisher.Id = user.Id
	media.Scope = entities.Public

	return media
}

func getBadMedia() *entities.Media {
	media := entities.MediaNew()
	media.Title = "Another lolcat video"
	media.Summary = "In this video, a lolcat play with catsnip !"
	user := getFooUser()
	media.Publisher.Name = user.Name
	media.Publisher.Email = user.Email
	media.Publisher.Id = user.Id
	media.Scope = entities.Private

	return media
}

func getFooPublicMedia() *entities.Media {
	media := entities.MediaNew()
	media.Title = "FooPublic"
	user := getFooUser()
	media.Publisher.Name = user.Name
	media.Publisher.Email = user.Email
	media.Publisher.Id = user.Id
	media.Scope = entities.Public

	return media
}

func getFooPrivateMedia() *entities.Media {
	media := entities.MediaNew()
	media.Title = "FooPrivate"
	user := getFooUser()
	media.Publisher.Name = user.Name
	media.Publisher.Email = user.Email
	media.Publisher.Id = user.Id
	media.Scope = entities.Private

	return media
}

func getFooLinkMedia() *entities.Media {
	media := entities.MediaNew()
	media.Title = "FooLink"
	user := getFooUser()
	media.Publisher.Name = user.Name
	media.Publisher.Email = user.Email
	media.Publisher.Id = user.Id
	media.Scope = entities.Link

	return media
}

func getAdminPublicMedia() *entities.Media {
	media := entities.MediaNew()
	media.Title = "AdminPublic"
	user := getAdminUser()
	media.Publisher.Name = user.Name
	media.Publisher.Email = user.Email
	media.Publisher.Id = user.Id
	media.Scope = entities.Public

	return media
}

func getAdminPrivateMedia() *entities.Media {
	media := entities.MediaNew()
	media.Title = "AdminPrivate"
	user := getAdminUser()
	media.Publisher.Name = user.Name
	media.Publisher.Email = user.Email
	media.Publisher.Id = user.Id
	media.Scope = entities.Private

	return media
}

func getAdminLinkMedia() *entities.Media {
	media := entities.MediaNew()
	media.Title = "AdminLink"
	user := getAdminUser()
	media.Publisher.Name = user.Name
	media.Publisher.Email = user.Email
	media.Publisher.Id = user.Id
	media.Scope = entities.Link

	return media
}
