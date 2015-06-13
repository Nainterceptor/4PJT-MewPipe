package users

import (
	"net/http"
	"supinfo/mewpipe/entities"
	"supinfo/mewpipe/filters"

	"github.com/Nainterceptor/go-restful"
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
		Returns(http.StatusOK, "User has been created", nil).
		Returns(http.StatusBadRequest, "Can't read entity", nil).
		Returns(http.StatusNotAcceptable, "Validation has failed", nil).
		Returns(http.StatusInternalServerError, "Return of MongoDB Insert", nil).
		Reads(entities.User{}))

	service.Route(service.
		GET("").
		Filter(filters.MustBeLogged).
		To(usersGet).
		Doc("Get a user list").
		Operation("usersGet").
		Param(service.QueryParameter("limit", "Limit number of results (default 25)").DataType("string")).
		Param(service.QueryParameter("start", "Start from item number n (default 0)").DataType("string")).
		Returns(http.StatusOK, "Users has been returned", nil).
		Returns(http.StatusInternalServerError, "Return of MongoDB find", nil))

	service.Route(service.
		PUT("/{user-id}").
		Filter(filters.MustBeLogged).
		Filter(filters.MustBeMyselfOrAdmin).
		To(userUpdate).
		Doc("Update a user").
		Operation("userUpdate").
		Param(service.PathParameter("user-id", "identifier of the user").DataType("string")).
		Returns(http.StatusOK, "User has been updated", nil).
		Returns(http.StatusBadRequest, "Bad Object ID or Can't read entity", nil).
		Returns(http.StatusNotFound, "User not found", nil).
		Returns(http.StatusNotAcceptable, "Validation has failed", nil).
		Returns(http.StatusInternalServerError, "Return of MongoDB Update", nil).
		Reads(entities.User{})) // from the request

	service.Route(service.
		PUT("/me").
		Filter(filters.MustBeLogged).
		Filter(filters.MustBeMyselfOrAdmin).
		To(userMeUpdate).
		Doc("Update my user").
		Operation("userMeUpdate").
		Returns(http.StatusOK, "User has been updated", nil).
		Returns(http.StatusBadRequest, "Bad Object ID or Can't read entity", nil).
		Returns(http.StatusNotFound, "User not found", nil).
		Returns(http.StatusNotAcceptable, "Validation has failed", nil).
		Returns(http.StatusInternalServerError, "Return of MongoDB Update", nil).
		Reads(entities.User{})) // from the request

	service.Route(service.
		DELETE("/{user-id}").
		Filter(filters.MustBeLogged).
		Filter(filters.MustBeMyselfOrAdmin).
		To(userDelete).
		Doc("Delete a user").
		Operation("userDelete").
		Param(service.PathParameter("user-id", "identifier of the user").DataType("string")).
		Returns(http.StatusNoContent, "User has been deleted", nil).
		Returns(http.StatusBadRequest, "Bad Object ID", nil).
		Returns(http.StatusNotFound, "User not found, eventually another MongoDB Fail", nil).
		Reads(entities.User{}))

	service.Route(service.
		DELETE("/me").
		Filter(filters.MustBeLogged).
		Filter(filters.MustBeMyselfOrAdmin).
		To(userMeDelete).
		Doc("Delete my user").
		Operation("userMeDelete").
		Returns(http.StatusNoContent, "User has been deleted", nil).
		Returns(http.StatusBadRequest, "Bad Object ID", nil).
		Returns(http.StatusNotFound, "User not found, eventually another MongoDB Fail", nil).
		Reads(entities.User{}))

	service.Route(service.
		GET("/{user-id}").
		To(userGet).
		Doc("Get a user").
		Operation("userGet").
		Param(service.PathParameter("user-id", "identifier of the user").DataType("string")).
		Returns(http.StatusOK, "User must be returned in the body", nil).
		Returns(http.StatusBadRequest, "Bad Object ID", nil).
		Returns(http.StatusNotFound, "User not found, eventually another MongoDB Fail", nil))

	service.Route(service.
		GET("/me").
		Filter(filters.MustBeLogged).
		Filter(filters.MustBeMyselfOrAdmin).
		To(userMeGet).
		Doc("Get my user").
		Operation("userMeGet").
		Returns(http.StatusOK, "User must be returned in the body", nil).
		Returns(http.StatusBadRequest, "Bad Object ID", nil).
		Returns(http.StatusNotFound, "User not found, eventually another MongoDB Fail", nil))

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

	service.Route(service.
		POST("/refresh-token").
		Filter(filters.MustBeLogged).
		To(userRefreshToken).
		Doc("Refresh a token").
		Notes("Get a new token from another token").
		Operation("userRefreshToken").
		Returns(http.StatusNotFound, "User not found, eventually another MongoDB Fail", nil).
		Returns(http.StatusInternalServerError, "Something failed while token generation", nil))

	container.Add(service)
}
