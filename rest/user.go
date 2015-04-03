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

    service.Route(service.POST("").To(CreateUser))
    return service
}

func CreateUser(request *restful.Request, response *restful.Response) {
    usr := entities.User{}
    err := request.ReadEntity(&usr)
    // here you would create the user with some persistence system
    if err == nil {
        response.WriteEntity(usr)
    } else {
        response.WriteError(http.StatusInternalServerError,err)
    }
}