package users

import (
	"net/http"

	"github.com/Nainterceptor/4PJT-MewPipe/entities"

	"github.com/emicklei/go-restful"
)

func userMeUpdate(request *restful.Request, response *restful.Response) {

	usr := request.Attribute("user").(*entities.User)

	oldUsr := *usr
	if err := request.ReadEntity(&usr); err != nil {
		response.WriteError(http.StatusBadRequest, err)
		return
	}

	//Persist some informations
	usr.Id = oldUsr.Id
	usr.Roles = oldUsr.Roles

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

func userMeDelete(request *restful.Request, response *restful.Response) {

	usr := request.Attribute("user").(*entities.User)

	if err := usr.Delete(); err != nil {
		response.WriteError(http.StatusNotFound, err)
		return
	}

	response.WriteHeader(http.StatusNoContent)
}

func userMeGet(request *restful.Request, response *restful.Response) {

	usr := request.Attribute("user").(*entities.User)
	response.WriteEntity(usr)
}
