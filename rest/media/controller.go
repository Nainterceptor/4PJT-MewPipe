package media

import (
	"net/http"
	"regexp"
	"strconv"
	"supinfo/mewpipe/entities"

	"github.com/emicklei/go-restful"
	"gopkg.in/mgo.v2/bson"
)

const (
	MAX_FILE_SIZE = 500 * 1000 * 1000
)

func mediaCreate(request *restful.Request, response *restful.Response) {
	user := request.Attribute("user").(*entities.User)

	media := entities.MediaNew()

	if err := request.ReadEntity(&media); err != nil {
		response.WriteError(http.StatusBadRequest, err)
		return
	}

	media.Publisher.Name = user.Name
	media.Publisher.Id = user.Id
	media.Publisher.Email = user.Email

	if err := media.Insert(); err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	response.WriteEntity(media)
}

func mediaUpload(request *restful.Request, response *restful.Response) {
	request.Request.Body = http.MaxBytesReader(response.ResponseWriter, request.Request.Body, MAX_FILE_SIZE)
	media := request.Attribute("media").(*entities.Media)

	request.Request.ParseMultipartForm(MAX_FILE_SIZE)
	mediaFile, handler, err := request.Request.FormFile("file")
	if err == nil {
		defer mediaFile.Close()
		go media.Upload(mediaFile, handler)

	}

	thumbnailFile, handler, err := request.Request.FormFile("thumbnail")
	if err == nil {
		defer thumbnailFile.Close()
		go media.UploadThumbnail(thumbnailFile, handler)
	}

	response.WriteEntity(media)
}

func mediaThumbnail(request *restful.Request, response *restful.Response) {
	media := request.Attribute("media").(*entities.Media)
	if media.Thumbnail == "" {
		response.WriteErrorString(http.StatusNotFound, "No thumbnail")
		return
	}
	if err := media.OpenThumbnail(); err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	defer media.CloseThumbnail()
	response.AddHeader("Content-type", "image/jpeg")
	response.AddHeader("Content-Length", strconv.FormatInt(media.ThumbnailSize(), 10))
	if err := media.CopyThumbnailTo(response); err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
}

func mediaRead(request *restful.Request, response *restful.Response) {
	media := request.Attribute("media").(*entities.Media)
	if media.File == "" {
		response.WriteErrorString(http.StatusNotFound, "No video")
		return
	}
	if err := media.OpenFile(); err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	defer media.CloseFile()

	response.AddHeader("Accept-Ranges", "bytes")
	response.AddHeader("Content-Disposition", "attachment; filename=video.mp4")
	response.AddHeader("Content-type", media.ContentType())
	if rangeReq := request.Request.Header.Get("range"); rangeReq != "" {
		regex, _ := regexp.Compile(`bytes=([0-9]*)-([0-9]*)`)
		ranges := regex.FindStringSubmatch(rangeReq)
		start := 0
		intSize := int(media.Size())
		end := intSize - 1
		if len(ranges) > 2 {
			testedStart, errStart := strconv.Atoi(ranges[1])
			if errStart == nil {
				start = testedStart
			}
			testedEnd, errEnd := strconv.Atoi(ranges[2])
			if errEnd == nil {
				end = testedEnd
			}

			if errStart == nil && errEnd != nil {
				start = testedStart
				end = intSize - 1
			} else if errStart != nil && errEnd == nil {
				start = intSize - testedEnd
				end = intSize - 1
			}
			err := media.SeekSet(int64(start))
			if err != nil {
				response.WriteError(http.StatusInternalServerError, err)
				return
			}
			currentSize := end + 1 - start
			buffer := make([]byte, currentSize)

			if err := media.Read(buffer); err != nil {
				response.WriteError(http.StatusInternalServerError, err)
				return
			}
			response.AddHeader("Content-Range", "bytes "+strconv.Itoa(start)+"-"+strconv.Itoa(end)+"/"+strconv.Itoa(intSize))
			response.AddHeader("Content-Length", strconv.FormatInt(int64(currentSize), 10))
			response.WriteHeader(http.StatusPartialContent)
			if _, err := response.ResponseWriter.Write(buffer); err != nil {
				response.WriteError(http.StatusInternalServerError, err)
				return
			}
			return
		}
	}
	response.AddHeader("Content-Length", strconv.FormatInt(media.Size(), 10))
	if err := media.CopyTo(response); err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
}

func mediaPut(request *restful.Request, response *restful.Response) {
	media := request.Attribute("media").(*entities.Media)
	oldMedia := *media

	if err := request.ReadEntity(&media); err != nil {
		response.WriteError(http.StatusBadRequest, err)
		return
	}

	//Prevent some editions
	media.Publisher = oldMedia.Publisher
	media.Id = oldMedia.Id
	media.File = oldMedia.File

	if err := media.Update(); err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	response.WriteEntity(media)
}

func mediaDelete(request *restful.Request, response *restful.Response) {
	media := request.Attribute("media").(*entities.Media)

	if err := media.Delete(); err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	response.WriteHeader(http.StatusNoContent)
}

func mediaGetAll(request *restful.Request, response *restful.Response) {
	user := request.Attribute("user")

	start, err := strconv.Atoi(request.QueryParameter("start"))
	if err != nil {
		start = 0
	}
	limit, err := strconv.Atoi(request.QueryParameter("limit"))
	if err != nil {
		limit = 25
	}
	if limit < 1 || limit > 100 {
		limit = 25
	}
	search := bson.M{}

	userParam := request.QueryParameter("user")
	if userParam != "" {
		if !bson.IsObjectIdHex(userParam) {
			response.WriteErrorString(http.StatusBadRequest, "Bad user Object ID")
			return
		}
		search["publisher._id"] = bson.ObjectIdHex(userParam)
	}

	orderParam := request.QueryParameter("order")
	order := "_id"

	regexOrder, _ := regexp.Compile("-?(_id|createdAt|views)")
	if regexOrder.MatchString(orderParam) {
		order = orderParam
	}

	var scopes []string
	scopes = append(scopes, "public")

	if user != nil {
		user := user.(*entities.User)
		if user.Id != "" {
			scopes = append(scopes, "private")
		}
	}
	search["scope"] = bson.M{"$in": scopes}

	medias, err := entities.MediaList(search, start, limit, order)

	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	response.WriteEntity(medias)
}

func mediaGet(request *restful.Request, response *restful.Response) {
	media := request.Attribute("media").(*entities.Media)
	userAttribute := request.Attribute("user")
	if userAttribute == nil {
		entities.ViewNewAnonymous(media.Id)
	} else {
		user := userAttribute.(*entities.User)
		entities.ViewNew(user.Id, media.Id)
	}
	media.CountViews()

	response.WriteEntity(media)
}

func mediaPostShare(request *restful.Request, response *restful.Response) {
	media := request.Attribute("media").(*entities.Media)
	userAttribute := request.Attribute("user")
	if userAttribute == nil {
		entities.ShareCountNewAnonymous(media.Id)
	} else {
		user := userAttribute.(*entities.User)
		entities.ShareCountNew(user.Id, media.Id)
	}
	media.CountShares()

	response.WriteEntity(media)
}
