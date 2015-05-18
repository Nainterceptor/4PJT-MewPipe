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
	if err := configs.MongoDB.C("media.chunks").DropCollection(); err != nil {
		fmt.Println(err)
	}
	if err := configs.MongoDB.C("media.files").DropCollection(); err != nil {
		fmt.Println(err)
	}
	if err := configs.MongoDB.C("media").DropCollection(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("All media deleted")
}

func getAmazingMedia() *entities.Media {
	media := entities.MediaNew()
	media.Title = "Amazing video"
	media.Summary = "In this video, a big comet hit Earth, and Superman dances with a unicorn."
	user := getFooUser()
	media.Publisher.Name = user.Name
	media.Publisher.Email = user.Email
	media.Publisher.Id = user.Id

	path, err := download("https://raw.githubusercontent.com/youtube/api-samples/master/java/src/main/resources/sample-video.mp4")
	if err == nil {
		in, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		header := textproto.MIMEHeader{
			"Content-Type": {"video/mp4; charset=UTF-8"},
		}
		fileHeader := new(multipart.FileHeader)
		fileHeader.Filename = "sample.mp4"
		fileHeader.Header = header
		media.Upload(in, fileHeader)
	}

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

	path, err := download("https://webglsamples.googlecode.com/hg-history/c4129c21f4d99e4d2138512c5a74554bdabb2257/color-adjust/sample-video.mp4")
	if err == nil {
		in, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		header := textproto.MIMEHeader{
			"Content-Type": {"video/mp4; charset=UTF-8"},
		}
		fileHeader := new(multipart.FileHeader)
		fileHeader.Filename = "sample.mp4"
		fileHeader.Header = header
		media.Upload(in, fileHeader)
	}

	return media
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
	}

	fmt.Println("Media added")
}
