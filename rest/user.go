package rest

import (
	"github.com/emicklei/go-restful"
	"supinfo/mewpipe/entities"
	"net/http"
    "gopkg.in/mgo.v2/bson"
    "time"
    "encoding/base64"
    "supinfo/mewpipe/filters"
)

func UserRoute() *restful.WebService {
	service := new(restful.WebService)
	service.
	Path("/rest/users").
	Consumes(restful.MIME_JSON).
	Produces(restful.MIME_JSON)

	service.Route(
        service.
            POST("").
            To(userCreate))
	service.Route(
        service.
            PUT("/{user-id}").
            Filter(filters.MustBeLogged).
            Filter(filters.UserIDMustBeMyself).
            To(userUpdate))
	service.Route(
        service.
            DELETE("/{user-id}").
            Filter(filters.MustBeLogged).
            Filter(filters.UserIDMustBeMyself).
            To(userDelete))
	service.Route(
        service.
            GET("/{user-id}").
            Filter(filters.MustBeLogged).
            Filter(filters.UserIDMustBeMyself).
            To(userGet))
	service.Route(
        service.
            POST("/login").
            To(userLogin))

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
    if err != nil {
        response.WriteError(http.StatusNotFound, err)
        return
    }
    if err := request.ReadEntity(&usr); err != nil {
        response.WriteError(http.StatusBadRequest, err)
        return
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

func userDelete(request *restful.Request, response *restful.Response) {

    id := request.PathParameter("user-id")
    if !bson.IsObjectIdHex(id) {
        response.WriteErrorString(http.StatusBadRequest, "Path must contain an Object ID")
        return
    }

    //userNew because find query is useless
    usr := entities.UserNewFromId(bson.ObjectIdHex(id))

    if err := usr.Delete(); err != nil {
        response.WriteError(http.StatusNotFound, err)
        return
    }

    response.WriteHeader(http.StatusNoContent)
}

func userGet(request *restful.Request, response *restful.Response) {

    id := request.PathParameter("user-id")
    if !bson.IsObjectIdHex(id) {
        response.WriteErrorString(http.StatusBadRequest, "Path must contain an Object ID")
        return
    }

    usr, err := entities.UserFromId(bson.ObjectIdHex(id))
    if err != nil {
        response.WriteError(http.StatusNotFound, err)
        return
    }

    response.WriteEntity(usr)
}

func userLogin(request *restful.Request, response *restful.Response) {

    form := entities.UserNew()

    if err := request.ReadEntity(&form); err != nil {
        response.WriteError(http.StatusBadRequest, err)
        return
    }
    usr, err := entities.UserFromCredentials(form.Email, form.Password)
    if err != nil {

    }
    token, err := usr.TokenNew()
    if err != nil {
        response.WriteError(http.StatusInternalServerError, err)
    }

    type lambdaReturn struct {
        User entities.User
        Token string
        ExpireAt time.Time
    }

    toReturn := new(lambdaReturn)
    toReturn.User = *usr
    toReturn.ExpireAt = token.ExpireAt
    toReturn.Token = base64.StdEncoding.EncodeToString([]byte(token.Token))

    response.WriteEntity(toReturn)
}