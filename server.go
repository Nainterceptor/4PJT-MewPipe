package main

import (
	"net/http"
	"supinfo/mewpipe/configs"
	"github.com/emicklei/go-restful"
	"supinfo/mewpipe/rest"
)

func main() {
	restful.DefaultContainer.Router(restful.CurlyRouter{})
	restful.Add(rest.UserRoute())
	restful.Add(configs.StaticRouter())

	http.ListenAndServe(":8080", nil)
}
