package configs
import (
    "fmt"
    "net/http"
    "path"

    "github.com/emicklei/go-restful"
)
func Router() *restful.WebService {
    ws := new(restful.WebService)

    ws.Route(ws.GET("/static/{subpath:*}").To(staticHandler))

    return ws
}

func staticHandler(req *restful.Request, resp *restful.Response) {
    actual := path.Join("static", req.PathParameter("subpath"))
    fmt.Printf("serving %s ... (from %s)\n", actual, req.PathParameter("subpath"))
    http.ServeFile(
    resp.ResponseWriter,
    req.Request,
    actual)
}