package configs
import (
    "fmt"
    "net/http"
    "path"

    "github.com/emicklei/go-restful"
)
func StaticRouter() *restful.WebService {
    ws := new(restful.WebService)

    ws.Route(ws.GET("/").To(homeHandler))
    ws.Route(ws.GET("/{subpath:*}").To(staticHandler))

    return ws
}

func staticHandler(req *restful.Request, resp *restful.Response) {
    actual := path.Join(*staticPath, req.PathParameter("subpath"))
    fmt.Printf("serving %s ... (from %s)\n", actual, req.PathParameter("subpath"))
    http.ServeFile(
    resp.ResponseWriter,
    req.Request,
    actual)
}

func homeHandler(req *restful.Request, resp *restful.Response) {
    http.ServeFile(resp.ResponseWriter, req.Request, path.Join(*staticPath, "index.html"))
}