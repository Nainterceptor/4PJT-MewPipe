package configs

import (
	"fmt"
	"net/http"
	"path"

	"os"

	"github.com/emicklei/go-restful"
)

func StaticRouter(container *restful.Container) {
	ws := new(restful.WebService)

	ws.Route(ws.GET("/").To(homeHandler))
	ws.Route(ws.GET("/{subpath:*}").To(staticHandler))
	container.Add(ws)
}

func staticHandler(req *restful.Request, resp *restful.Response) {
	actual := path.Join(*staticPath, req.PathParameter("subpath"))
	if _, err := os.Stat(actual); os.IsNotExist(err) {
		actual = path.Join(*staticPath, "index.html")
	}
	fmt.Printf("serving %s ... (from %s)\n", actual, req.PathParameter("subpath"))
	http.ServeFile(
		resp.ResponseWriter,
		req.Request,
		actual)
}

func homeHandler(req *restful.Request, resp *restful.Response) {
	http.ServeFile(resp.ResponseWriter, req.Request, path.Join(*staticPath, "index.html"))
}
