package media

import (
	"net/http"
	"regexp"
	"strconv"
	"supinfo/mewpipe/entities"

	"github.com/emicklei/go-restful"
)

func mediaCreate(request *restful.Request, response *restful.Response) {

	media := entities.MediaNew()

	if err := request.ReadEntity(&media); err != nil {
		response.WriteError(http.StatusBadRequest, err)
		return
	}

	if err := media.Insert(); err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	response.WriteEntity(media)
}

func mediaUpload(request *restful.Request, response *restful.Response) {
	media := request.Attribute("media").(*entities.Media)

	request.Request.ParseMultipartForm(500 * 1000 * 1000)
	postedFile, handler, err := request.Request.FormFile("file")
	if err != nil {
		response.WriteError(http.StatusBadRequest, err)
		return
	}
	defer postedFile.Close()

	if err := media.Upload(postedFile, handler); err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	response.WriteEntity(media)
}

func mediaRead(request *restful.Request, response *restful.Response) {
	media := request.Attribute("media").(*entities.Media)
	if err := media.OpenFile(); err != nil {

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

	if err := request.ReadEntity(&media); err != nil {
		response.WriteError(http.StatusBadRequest, err)
		return
	}

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

func mediaGet(request *restful.Request, response *restful.Response) {
	media := request.Attribute("media").(*entities.Media)

	response.WriteEntity(media)
}
