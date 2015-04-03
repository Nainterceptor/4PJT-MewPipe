package rest
import (
    "github.com/emicklei/go-restful"
    "supinfo/mewpipe/entities"
    "net/http"
    "log"
    "gopkg.in/mgo.v2/bson"
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
    errRE := request.ReadEntity(&usr)
    usr.Id = bson.NewObjectId()
    err := entities.UserCollection.Insert(&usr)
    if err != nil {
        log.Fatal(err)
    }

    // here you would create the user with some persistence system
    if errRE == nil {
        response.WriteEntity(usr)
    } else {
        response.WriteError(http.StatusInternalServerError,errRE)
    }
}