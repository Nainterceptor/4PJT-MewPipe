package filters

import (
	"net/http"
	"supinfo/mewpipe/entities"

	"github.com/emicklei/go-restful"
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
		resp.WriteErrorString(http.StatusNotFound, "Media not found")
		return
	}
	req.SetAttribute("media", media)
	chain.ProcessFilter(req, resp)
}

func MustBeOwnerOrAdmin(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	usr := request.Attribute("user").(*entities.User)
	media := request.Attribute("media").(*entities.Media)

	if media.Publisher.Id != usr.Id && !usr.HasRole("Admin") {
		response.WriteErrorString(http.StatusForbidden, "You must be owner or admin")
		return
	}
	chain.ProcessFilter(request, response)
}

func ScopeControl(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	user := request.Attribute("user")
	media := request.Attribute("media").(*entities.Media)
	if media.Scope == "private" && (user == nil || user.(*entities.User).Id == "") {
		response.WriteErrorString(http.StatusForbidden, "This video is private")
		return
	}
	chain.ProcessFilter(request, response)
}
