package rest

import (
    "github.com/emicklei/go-restful"
    "supinfo/mewpipe/entities"
    "net/http"
    "io"
)

func MediaRoute() *restful.WebService {
    service := new(restful.WebService)
    service.
    Path("/rest/media").
    Consumes(restful.MIME_JSON).
    Produces(restful.MIME_JSON)

    service.Route(service.POST("").To(mediaCreate))
    service.Route(service.POST("/upload").Consumes("multipart/form-data").To(mediaUpload))

    return service
}

func mediaCreate(request *restful.Request, response *restful.Response) {

    media := entities.Media{}
    errRE := request.ReadEntity(&media)

    if err := entities.MediaCollection.Insert(&media); err != nil {
        response.WriteError(http.StatusInternalServerError, err)
        return
    }

    if errRE == nil {
        response.WriteEntity(media)
    } else {
        response.WriteError(http.StatusInternalServerError, errRE)
    }
}

func mediaUpload(request *restful.Request, response *restful.Response) {
    request.Request.ParseMultipartForm(500 * 1000 * 1000)
    postedFile, handler, err := request.Request.FormFile("file")
    if err != nil {
        response.WriteError(http.StatusInternalServerError, err)
    }
    defer postedFile.Close()
    mongoFile, err := entities.MediaGridFS.Create(handler.Filename)
    defer mongoFile.Close()
    io.Copy(mongoFile, postedFile)
    if err != nil {
        response.WriteError(http.StatusInternalServerError, err)
    }
}
