package entities

import (
	"fmt"
	"mime/multipart"
	"net/textproto"
	"os"
	"supinfo/mewpipe/configs"
	"supinfo/mewpipe/entities"
)

func ClearMedia() {
	configs.MongoDB.C("media.chunks").DropCollection()
	configs.MongoDB.C("media.thumbnails.chunks").DropCollection()
	configs.MongoDB.C("media.files").DropCollection()
	configs.MongoDB.C("media.thumbnails.files").DropCollection()
	configs.MongoDB.C("media").DropCollection()
	fmt.Println("All media deleted")
}

func InsertSomeMedia() {
	var mediaArray []*entities.Media
	mediaArray = append(mediaArray, getBadMedia())
	mediaArray = append(mediaArray, getAmazingMedia())
	for _, media := range mediaArray {
		if err := media.Insert(); err != nil {
			panic(err)
			return
		}
		addVideoBin(media)

	}
	addVideoBin(mediaArray[1])
	fmt.Println("Media added")
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
