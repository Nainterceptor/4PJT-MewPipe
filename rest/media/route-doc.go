package media

import (
	"net/http"
	"supinfo/mewpipe/entities"
	"supinfo/mewpipe/filters"

	"github.com/emicklei/go-restful"
)

func MediaRoute(container *restful.Container) {
	service := new(restful.WebService)
	service.
		Path("/rest/media").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	service.Route(service.
		POST("").
		Filter(filters.MustBeLogged).
		To(mediaCreate).
		Doc("New vidéo (metadocument)").
		Operation("mediaCreate").
		Returns(http.StatusOK, "User, token, and ExpirationDate are returned", nil).
		Returns(http.StatusBadRequest, "Can't read entity", nil).
		Returns(http.StatusInternalServerError, "MongoDB fail while insertion", nil).
		Reads(entities.Media{}))

	service.Route(service.
		PUT("/{media-id}").
		Filter(filters.InjectMediaMeta).
		Filter(filters.MustBeLogged).
		Filter(filters.MustBeOwnerOrAdmin).
		To(mediaPut).
		Doc("Update a vidéo (metadocument)").
		Operation("mediaPut").
		Param(service.PathParameter("media-id", "identifier of the media").DataType("string")).
		Returns(http.StatusOK, "Video has been updated", nil).
		Returns(http.StatusBadRequest, "Bad ID / Can't read entity", nil).
		Returns(http.StatusInternalServerError, "MongoDB fail while update", nil).
		Reads(entities.Media{}))

	service.Route(service.
		GET("").
		Filter(filters.InjectUser).
		To(mediaGetAll).
		Doc("Get a media list").
		Operation("mediaGetAll").
		Param(service.QueryParameter("user", "Filter user media").DataType("string")).
		Param(service.QueryParameter("order", "Order by ... ?").DataType("string")).
		Param(service.QueryParameter("limit", "Limit number of results (default 25)").DataType("string")).
		Param(service.QueryParameter("start", "Start from item number n (default 0)").DataType("string")).
		Returns(http.StatusOK, "Medias has been returned", nil).
		Returns(http.StatusInternalServerError, "Return of MongoDB find", nil))

	service.Route(service.
		GET("/{media-id}").
		Filter(filters.InjectMediaMeta).
		Filter(filters.InjectUser).
		Filter(filters.ScopeControl).
		To(mediaGet).
		Doc("Get a vidéo (metadocument)").
		Operation("mediaGet").
		Param(service.PathParameter("media-id", "identifier of the media").DataType("string")).
		Returns(http.StatusOK, "Video has been returned", nil).
		Returns(http.StatusBadRequest, "Bad ID", nil).
		Returns(http.StatusNotFound, "Media not found", nil))

	service.Route(service.
		POST("/{media-id}/share").
		Filter(filters.InjectMediaMeta).
		Filter(filters.InjectUser).
		Filter(filters.ScopeControl).
		To(mediaPostShare).
		Doc("Set a new share and return media").
		Operation("mediaPostShare").
		Param(service.PathParameter("media-id", "identifier of the media").DataType("string")).
		Returns(http.StatusOK, "Video has been returned", nil).
		Returns(http.StatusBadRequest, "Bad ID", nil).
		Returns(http.StatusNotFound, "Media not found", nil))

	service.Route(service.
		DELETE("/{media-id}").
		Filter(filters.InjectMediaMeta).
		Filter(filters.MustBeLogged).
		Filter(filters.MustBeOwnerOrAdmin).
		To(mediaDelete).
		Doc("Delete a vidéo (metadocument & bin)").
		Operation("mediaDelete").
		Param(service.PathParameter("media-id", "identifier of the media").DataType("string")).
		Returns(http.StatusNoContent, "Video has been deleted", nil).
		Returns(http.StatusBadRequest, "Bad ID", nil).
		Returns(http.StatusNotFound, "Media not found", nil).
		Returns(http.StatusInternalServerError, "MongoDB fail while delete", nil))

	service.Route(service.
		POST("/{media-id}/upload").
		Filter(filters.InjectMediaMeta).
		Filter(filters.MustBeLogged).
		Filter(filters.MustBeOwnerOrAdmin).
		Consumes("multipart/form-data").
		To(mediaUpload).
		Doc("Upload a vidéo (bin)").
		Operation("mediaUpload").
		Param(service.PathParameter("media-id", "identifier of the media").DataType("string")).
		Returns(http.StatusOK, "Video has been uploaded", nil).
		Returns(http.StatusBadRequest, "Bad ID or bad parsing", nil).
		Returns(http.StatusNotFound, "Media not found", nil).
		Returns(http.StatusInternalServerError, "MongoDB fail while uploading", nil))

	service.Route(service.
		GET("/{media-id}/read").
		Filter(filters.InjectMediaMeta).
		Filter(filters.InjectUser).
		Filter(filters.ScopeControl).
		To(mediaRead).
		Doc("Play a vidéo (bin)").
		Notes("Can handle a range request").
		Operation("mediaRead").
		Param(service.PathParameter("media-id", "identifier of the media").DataType("string")).
		Returns(http.StatusOK, "Video has been read", nil).
		Returns(http.StatusPartialContent, "Video part (range request)", nil).
		Returns(http.StatusBadRequest, "Bad ID", nil).
		Returns(http.StatusNotFound, "Media not found", nil).
		Returns(http.StatusInternalServerError, "MongoDB fail while reading (Seek, Read, Write)", nil))

	service.Route(service.
		GET("/{media-id}/thumbnail").
		Filter(filters.InjectMediaMeta).
		Filter(filters.InjectUser).
		Filter(filters.ScopeControl).
		To(mediaThumbnail).
		Doc("Get a thumbnail (bin)").
		Operation("mediaThumbnail").
		Param(service.PathParameter("media-id", "identifier of the media").DataType("string")).
		Returns(http.StatusOK, "Image has been read", nil).
		Returns(http.StatusBadRequest, "Bad ID", nil).
		Returns(http.StatusNotFound, "Media not found", nil).
		Returns(http.StatusInternalServerError, "MongoDB fail while reading (Seek, Read, Write)", nil))

	container.Add(service)
}
