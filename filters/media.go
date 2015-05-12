package filters
import (
    "github.com/emicklei/go-restful"
    "supinfo/mewpipe/entities"
    "net/http"
    "gopkg.in/mgo.v2/bson"
)

func InjectMediaMeta(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
    id := req.PathParameter("media-id")
    if !bson.IsObjectIdHex(id) {
        resp.WriteErrorString(http.StatusBadRequest, "Path must contain an Object ID")
        return
    }
    oid := bson.ObjectIdHex(id)
    media, err := entities.MediaFromId(oid)
    if err != nil {
        resp.WriteErrorString(http.StatusForbidden, "Media not found")
        return
    }
    req.SetAttribute("media", media)
    chain.ProcessFilter(req, resp)
}