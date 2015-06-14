package filters

import (
	"net/http"
	"supinfo/mewpipe/entities"

	"github.com/emicklei/go-restful"
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

func InjectUser(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {

	token := req.Request.Header.Get("Authorization")
	if token == "" {
		cookie, err := req.Request.Cookie("accessToken")
		if err == nil {
			token = cookie.Value
		}
	}
	if token != "" {
		usr, err := entities.UserFromToken(token)
		if err == nil {
			req.SetAttribute("user", usr)
		}
	}

	chain.ProcessFilter(req, resp)
}

func MustBeMyselfOrAdmin(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {

	id := req.PathParameter("user-id")
	if !bson.IsObjectIdHex(id) {
		resp.WriteErrorString(http.StatusBadRequest, "Path must contain an Object ID")
		return
	}
	usr := req.Attribute("user").(*entities.User)
	if bson.ObjectIdHex(id) != usr.Id && !usr.HasRole("Admin") {
		resp.WriteErrorString(http.StatusForbidden, "You must be owner or admin")
		return
	}
	chain.ProcessFilter(req, resp)
}
