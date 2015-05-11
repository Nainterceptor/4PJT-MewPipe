package media

import (
    "github.com/emicklei/go-restful"
    "supinfo/mewpipe/entities"
    "net/http"
)

func MediaRoute(container *restful.Container) {
    service := new(restful.WebService)
    service.
        Path("/rest/media").
        Consumes(restful.MIME_JSON).
        Produces(restful.MIME_JSON)

    service.Route(service.
        POST("").
        To(mediaCreate).

        Doc("New vidéo (metadocument)").
        Operation("mediaCreate").
        Returns(http.StatusOK, "User, token, and ExpirationDate are returned", nil).
        Returns(http.StatusBadRequest, "Can't read entity", nil).
        Returns(http.StatusInternalServerError, "MongoDB fail while insertion", nil).
        Reads(entities.Media{}))

    service.Route(service.
        PUT("/{media-id}").
        To(mediaPut).

        Doc("Update a vidéo (metadocument)").
        Operation("mediaPut").
        Param(service.PathParameter("user-id", "identifier of the media").DataType("string")).
        Returns(http.StatusOK, "Video has been updated", nil).
        Returns(http.StatusBadRequest, "Bad ID / Can't read entity", nil).
        Returns(http.StatusInternalServerError, "MongoDB fail while update", nil).
        Reads(entities.Media{}))

    service.Route(service.
        GET("/{media-id}").
        To(mediaGet))

    service.Route(service.
        DELETE("/{media-id}").
        To(mediaDelete))

    service.Route(service.
        POST("/{media-id}/upload").
        Consumes("multipart/form-data").
        To(mediaUpload))

    service.Route(service.
        GET("/{media-id}/read").
        To(mediaRead))

    container.Add(service)
}