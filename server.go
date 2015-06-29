package main

import (
	"net/http"

	"github.com/Nainterceptor/4PJT-MewPipe/configs"
	"github.com/Nainterceptor/4PJT-MewPipe/rest/media"
	"github.com/Nainterceptor/4PJT-MewPipe/rest/users"

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
