package rest

import (
	"github.com/emicklei/go-restful"
	"supinfo/mewpipe/entities"
	"net/http"
    "gopkg.in/mgo.v2/bson"
)

func UserRoute() *restful.WebService {
	service := new(restful.WebService)
	service.
	Path("/rest/users").
	Consumes(restful.MIME_JSON).
	Produces(restful.MIME_JSON)

	service.Route(service.POST("").To(userCreate))
	service.Route(service.PUT("/{user-id}").To(userUpdate))

	return service
}

func userCreate(request *restful.Request, response *restful.Response) {

	usr := entities.UserNew();

	if err := request.ReadEntity(&usr); err != nil {
		response.WriteError(http.StatusBadRequest, err)
	}

    if err := usr.Validate(); err != nil {
        response.WriteError(http.StatusNotAcceptable, err)
        return
    }

	if err := usr.Insert(); err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

    response.WriteEntity(usr)
}

func userUpdate(request *restful.Request, response *restful.Response) {

    id := request.PathParameter("user-id")
    if !bson.IsObjectIdHex(id) {
        response.WriteErrorString(http.StatusBadRequest, "Path must contain an Object ID")
        return
    }

    usr, err := entities.UserFromId(bson.ObjectIdHex(id))
    if (err != nil) {
        response.WriteError(http.StatusNotFound, err)
    }
    if err := request.ReadEntity(&usr); err != nil {
        response.WriteError(http.StatusBadRequest, err)
    }

    if err := usr.Validate(); err != nil {
        response.WriteError(http.StatusNotAcceptable, err)
        return
    }

    if err := usr.Update(); err != nil {
        response.WriteError(http.StatusInternalServerError, err)
        return
    }

    response.WriteEntity(usr)
}