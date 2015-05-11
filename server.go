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
    restful.Add(configs.StaticRouter())

    configs.ConfigureSwagger(wsContainer)

    restful.DefaultContainer.Router(restful.CurlyRouter{})

    server := &http.Server{Addr: *configs.HttpBinding, Handler: wsContainer}
    server.ListenAndServe()
}
