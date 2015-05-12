package main

import (
	"net/http"
	"supinfo/mewpipe/configs"
	"github.com/emicklei/go-restful"
	"supinfo/mewpipe/rest/users"
	"supinfo/mewpipe/rest/media"
    "log"
)

func main() {
    configs.Parse()

    container := restful.NewContainer()
    configs.StaticRouter(container)

    go func() {
        server := &http.Server{Addr: ":8181", Handler: container}
        server.ListenAndServe()
    }()

    wsContainer := restful.NewContainer()
    users.UserRoute(wsContainer)
    media.MediaRoute(wsContainer)

    configs.ConfigureSwagger(wsContainer)
    server := &http.Server{Addr: *configs.HttpBinding, Handler: wsContainer}
    log.Fatal(server.ListenAndServe())
}
