package main

import (
	"net/http"
	"supinfo/mewpipe/configs"
	"github.com/emicklei/go-restful"
	"supinfo/mewpipe/rest"
)

func main() {
    configs.Parse()

    wsContainer := restful.NewContainer()
    rest.UserRoute(wsContainer)
    rest.MediaRoute(wsContainer)
    configs.StaticRouter(wsContainer)

    configs.ConfigureSwagger(wsContainer)


    server := &http.Server{Addr: *configs.HttpBinding, Handler: wsContainer}
    server.ListenAndServe()
}
