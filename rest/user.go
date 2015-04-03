package rest
import (
    "github.com/emicklei/go-restful"
    "supinfo/mewpipe/entities"
    "gopkg.in/mgo.v2"
    "net/http"
    "supinfo/mewpipe/configs"
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
    session, err := mgo.Dial(*configs.MongoCS)
    if err != nil {
        panic(err)
    }
    defer session.Close()

    usr := entities.User{}
    errRE := request.ReadEntity(&usr)
    // here you would create the user with some persistence system
    if errRE == nil {
        response.WriteEntity(usr)
    } else {
        response.WriteError(http.StatusInternalServerError,errRE)
    }
}