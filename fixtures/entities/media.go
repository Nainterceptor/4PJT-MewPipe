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
	media.Scope = "public"

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
	media.Scope = "private"

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
	addVideo(mediaArray[0], "https://raw.githubusercontent.com/youtube/api-samples/master/java/src/main/resources/sample-video.mp4")
	addVideo(mediaArray[1], "https://webglsamples.googlecode.com/hg-history/c4129c21f4d99e4d2138512c5a74554bdabb2257/color-adjust/sample-video.mp4")
	fmt.Println("Media added")
}

func addVideo(media *entities.Media, link string) {
	path, err := download(link)
	if err != nil {
		fmt.Println(err)
		return
	}
	in, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	header := textproto.MIMEHeader{
		"Content-Type": {"video/mp4; charset=UTF-8"},
	}
	fileHeader := new(multipart.FileHeader)
	fileHeader.Filename = "sample.mp4"
	fileHeader.Header = header
	if err = media.Upload(in, fileHeader); err != nil {
		fmt.Println(err)
		return
	}
}
