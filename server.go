package main

import (
	"net/http"
	"supinfo/mewpipe/configs"
	"github.com/emicklei/go-restful"
)

func main() {
	restful.DefaultContainer.Router(restful.CurlyRouter{})

	ws := configs.Router()
	restful.Add(ws)

	http.ListenAndServe(":8080", nil)
}
