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

func UserRoute(container *restful.Container) {
	service := new(restful.WebService)
	service.
        Path("/rest/users").
        Consumes(restful.MIME_JSON).
        Produces(restful.MIME_JSON)

	service.Route(service.
        POST("").
        To(userCreate).

        Doc("Create a new user").
        Operation("userCreate").
        Param(service.PathParameter("user-id", "identifier of the user").DataType("string")).
        Returns(http.StatusOK, "User has been created", nil).
        Returns(http.StatusBadRequest, "Can't read entity", nil).
        Returns(http.StatusNotAcceptable, "Validation has failed", nil).
        Returns(http.StatusInternalServerError, "Return of MongoDB Insert", nil).
        Reads(entities.User{}))

    service.Route(service.
        PUT("/{user-id}").
        Filter(filters.MustBeLogged).
        Filter(filters.UserIDMustBeMyself).
        To(userUpdate).

        Doc("Update a user").
        Operation("userUpdate").
        Param(service.PathParameter("user-id", "identifier of the user").DataType("string")).
        Returns(http.StatusOK, "User has been updated", nil).
        Returns(http.StatusBadRequest, "Bad Object ID or Can't read entity", nil).
        Returns(http.StatusNotFound, "User not found", nil).
        Returns(http.StatusNotAcceptable, "Validation has failed", nil).
        Returns(http.StatusInternalServerError, "Return of MongoDB Update", nil).
        Reads(entities.User{}))// from the request

	service.Route(service.
        DELETE("/{user-id}").
        Filter(filters.MustBeLogged).
        Filter(filters.UserIDMustBeMyself).
        To(userDelete).

        Doc("Delete a user").
        Operation("userDelete").
        Param(service.PathParameter("user-id", "identifier of the user").DataType("string")).
        Returns(http.StatusNoContent, "User has been deleted", nil).
        Returns(http.StatusBadRequest, "Bad Object ID", nil).
        Returns(http.StatusNotFound, "User not found, eventually another MongoDB Fail", nil).
        Reads(entities.User{}))

	service.Route(service.
        GET("/{user-id}").
        Filter(filters.MustBeLogged).
        Filter(filters.UserIDMustBeMyself).
        To(userGet).

        Doc("Get a user").
        Operation("userGet").
        Param(service.PathParameter("user-id", "identifier of the user").DataType("string")).
        Returns(http.StatusOK, "User must be returned in the body", nil).
        Returns(http.StatusBadRequest, "Bad Object ID", nil).
        Returns(http.StatusNotFound, "User not found, eventually another MongoDB Fail", nil).
        Reads(entities.User{}))

	service.Route(service.
        POST("/login").
        To(userLogin).

        Doc("Login").
        Notes("Only Password and Email are important").
        Operation("userLogin").
        Returns(http.StatusOK, "User, token, and ExpirationDate are returned", nil).
        Returns(http.StatusBadRequest, "Can't read entity", nil).
        Returns(http.StatusNotFound, "User not found, eventually another MongoDB Fail about authentication", nil).
        Returns(http.StatusInternalServerError, "Something failed while token generation", nil).
        Reads(entities.User{}))


    container.Add(service)
}

func userCreate(request *restful.Request, response *restful.Response) {

	usr := entities.UserNew();

	if err := request.ReadEntity(&usr); err != nil {
		response.WriteError(http.StatusBadRequest, err)
        return
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
        response.WriteError(http.StatusNotFound, err)
        return
    }
    token, err := usr.TokenNew()
    if err != nil {
        response.WriteError(http.StatusInternalServerError, err)
        return
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