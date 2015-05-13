package main

import (
	"net/http"
	"supinfo/mewpipe/configs"
	"github.com/emicklei/go-restful"
	"supinfo/mewpipe/rest/users"
	"supinfo/mewpipe/rest/media"
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
