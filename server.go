package main

import (
	"net/http"
	"supinfo/mewpipe/configs"
	"github.com/emicklei/go-restful"
	"supinfo/mewpipe/rest"
)

func main() {
	configs.Parse()
	restful.DefaultContainer.Router(restful.CurlyRouter{})
	restful.Add(rest.UserRoute())
	restful.Add(configs.StaticRouter())
	http.ListenAndServe(*configs.HttpBinding, nil)
}
