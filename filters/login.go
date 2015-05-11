package filters

import (
    "github.com/emicklei/go-restful"
    "supinfo/mewpipe/entities"
    "net/http"
    "gopkg.in/mgo.v2/bson"
)

func MustBeLogged(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {

    usr, err := entities.UserFromToken(req.Request.Header.Get("Authorization"))
    if err != nil {
        resp.WriteErrorString(http.StatusForbidden, "User not found")
        return
    }

    req.SetAttribute("user", usr)
    chain.ProcessFilter(req, resp)
}

func UserIDMustBeMyself(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {

    id := req.PathParameter("user-id")
    if !bson.IsObjectIdHex(id) {
        resp.WriteErrorString(http.StatusBadRequest, "Path must contain an Object ID")
        return
    }
    usr := req.Attribute("user").(*entities.User)
    if bson.ObjectIdHex(id) != usr.Id {
        resp.WriteErrorString(http.StatusForbidden, "Can't edit")
        return
    }
    chain.ProcessFilter(req, resp)
}