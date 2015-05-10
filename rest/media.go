package rest

import (
    "github.com/emicklei/go-restful"
    "supinfo/mewpipe/entities"
    "net/http"
    "io"
    "gopkg.in/mgo.v2/bson"
    "strconv"
    "regexp"
    "os"
)

func MediaRoute() *restful.WebService {
    service := new(restful.WebService)
    service.
    Path("/rest/media").
    Consumes(restful.MIME_JSON).
    Produces(restful.MIME_JSON)

    service.Route(service.POST("").To(mediaCreate))
    service.Route(service.PUT("/{media-id}").To(mediaPut))
    service.Route(service.POST("/{media-id}/upload").Consumes("multipart/form-data").To(mediaUpload))
    service.Route(service.GET("/{media-id}/read").To(mediaRead))

    return service
}

func mediaCreate(request *restful.Request, response *restful.Response) {

    media := entities.Media{}
    if err := request.ReadEntity(&media); err != nil {
        response.WriteError(http.StatusInternalServerError, err)
        return
    }
    media.Id = bson.NewObjectId()

    if err := entities.MediaCollection.Insert(&media); err != nil {
        response.WriteError(http.StatusInternalServerError, err)
        return
    }
    response.WriteEntity(media)
}

func mediaUpload(request *restful.Request, response *restful.Response) {
    id := request.PathParameter("media-id")
    if !bson.IsObjectIdHex(id) {
        response.WriteErrorString(404, "Bad ID")
        return
    }
    oid := bson.ObjectIdHex(id)
    media := entities.Media{}

    if err := entities.MediaCollection.FindId(oid).One(&media); err != nil {
        response.WriteError(http.StatusInternalServerError, err)
        return
    }

    request.Request.ParseMultipartForm(500 * 1000 * 1000)
    postedFile, handler, err := request.Request.FormFile("file")
    if err != nil {
        response.WriteError(http.StatusInternalServerError, err)
        return
    }
    defer postedFile.Close()
    mongoFile, err := entities.MediaGridFS.Create(handler.Filename)
    defer mongoFile.Close()
    if err != nil {
        response.WriteError(http.StatusInternalServerError, err)
        return
    }
    io.Copy(mongoFile, postedFile)
    mongoFile.SetContentType(handler.Header.Get("Content-Type"))
    media.File = mongoFile.Id().(bson.ObjectId)
    if err := entities.MediaCollection.UpdateId(oid,&media); err != nil {
        response.WriteError(http.StatusInternalServerError, err)
        mongoFile.Abort()
        return
    }
    response.WriteEntity(media)
}

func mediaRead(request *restful.Request, response *restful.Response) {
    id := request.PathParameter("media-id")
    if !bson.IsObjectIdHex(id) {
        response.WriteErrorString(404, "Bad ID")
        return
    }
    oid := bson.ObjectIdHex(id)
    media := entities.Media{}
    if err := entities.MediaCollection.FindId(oid).One(&media); err != nil {
        response.WriteError(http.StatusInternalServerError, err)
        return
    }

    mongoFile, err := entities.MediaGridFS.OpenId(media.File)
    defer mongoFile.Close()
    if err != nil {
        response.WriteError(http.StatusInternalServerError, err)
        return
    }
    response.AddHeader("Accept-Ranges", "bytes")
    response.AddHeader("Content-Disposition", "attachment; filename=video.mp4")
    response.AddHeader("Content-type", mongoFile.ContentType())
    if rangeReq := request.Request.Header.Get("range"); rangeReq != "" {
        regex, _ := regexp.Compile(`bytes=([0-9]*)-([0-9]*)`)
        ranges := regex.FindStringSubmatch(rangeReq)
        start := 0
        intSize := int(mongoFile.Size())
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
            _, err := mongoFile.Seek(int64(start), os.SEEK_SET)
            if err != nil {
                response.WriteError(http.StatusInternalServerError, err)
                return
            }
            currentSize := end + 1 - start
            buffer := make([]byte, currentSize)
            _, err = mongoFile.Read(buffer)
            if err != nil {
                response.WriteError(http.StatusInternalServerError, err)
                return
            }
            response.AddHeader("Content-Range", "bytes " + strconv.Itoa(start) + "-" + strconv.Itoa(end) + "/" + strconv.Itoa(intSize))
            response.AddHeader("Content-Length", strconv.FormatInt(int64(currentSize), 10))
            response.WriteHeader(http.StatusPartialContent)
            if _, err := response.ResponseWriter.Write(buffer); err != nil {
                response.WriteError(http.StatusInternalServerError, err)
                return
            }
            return
        }
    }
    response.AddHeader("Content-Length", strconv.FormatInt(mongoFile.Size(), 10))
    _, err = io.Copy(response, mongoFile)
}

func mediaPut(request *restful.Request, response *restful.Response) {
    id := request.PathParameter("media-id")
    if !bson.IsObjectIdHex(id) {
        response.WriteErrorString(400, "Bad ID")
        return
    }
    oid := bson.ObjectIdHex(id)
    media := entities.Media{}
    if err := entities.MediaCollection.FindId(oid).One(&media); err != nil {
        response.WriteError(http.StatusInternalServerError, err)
        return
    }

    if err := request.ReadEntity(&media); err != nil {
        response.WriteError(http.StatusInternalServerError, err)
        return
    }

    if _, err := entities.MediaCollection.UpsertId(oid, &media); err != nil {
        response.WriteError(http.StatusInternalServerError, err)
        return
    }

    response.WriteEntity(media)
}