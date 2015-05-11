package rest

import (
	"github.com/emicklei/go-restful"
	"supinfo/mewpipe/entities"
	"net/http"
)

func UserRoute() *restful.WebService {
	service := new(restful.WebService)
	service.
	Path("/rest/users").
	Consumes(restful.MIME_JSON).
	Produces(restful.MIME_JSON)

	service.Route(service.POST("").To(createUser))

	return service
}

func createUser(request *restful.Request, response *restful.Response) {

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