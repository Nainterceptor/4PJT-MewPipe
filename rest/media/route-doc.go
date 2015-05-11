package media

import (
    "github.com/emicklei/go-restful"
)

func MediaRoute(container *restful.Container) {
    service := new(restful.WebService)
    service.
    Path("/rest/media").
    Consumes(restful.MIME_JSON).
    Produces(restful.MIME_JSON)

    service.Route(service.POST("").To(mediaCreate))
    service.Route(service.PUT("/{media-id}").To(mediaPut))
    service.Route(service.GET("/{media-id}").To(mediaGet))
    service.Route(service.DELETE("/{media-id}").To(mediaDelete))
    service.Route(service.POST("/{media-id}/upload").Consumes("multipart/form-data").To(mediaUpload))
    service.Route(service.GET("/{media-id}/read").To(mediaRead))

    container.Add(service)
}