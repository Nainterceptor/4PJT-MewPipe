package main

import (
	"net/http"
	"supinfo/mewpipe/configs"
	"supinfo/mewpipe/rest/media"
	"supinfo/mewpipe/rest/users"

	"github.com/emicklei/go-restful"
)

func main() {

	wsContainer := restful.NewContainer()
	users.UserRoute(wsContainer)
	media.MediaRoute(wsContainer)
	configs.StaticRouter(wsContainer)

	configs.ConfigureSwagger(wsContainer)

	server := &http.Server{Addr: *configs.HttpBinding, Handler: wsContainer}
	server.ListenAndServe()
}
