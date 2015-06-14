package configs

import (
	"flag"

	"github.com/Nainterceptor/go-restful"
	"github.com/Nainterceptor/go-restful/swagger"
	"github.com/vharitonsky/iniflags"
	"gopkg.in/mgo.v2"
)

var staticPath = flag.String("static_path", "static", "Localisation for static files")

var HttpBinding = flag.String("http_binding", "localhost:1337", "IP/Port to listen HTTP Server")
var ServerDomain = flag.String("server_domain", "localhost:1337", "Useful to forward to OpenID")

var mongoCS = flag.String("mongodb_CS", "localhost", "Connection endpoint for mongodb driver")
var mongoName = flag.String("mongodb_DB", "MewPipe", "Database to mount")

var TWITTER_KEY = flag.String("TWITTER_KEY", "", "Twitter app key")
var TWITTER_SECRET = flag.String("TWITTER_SECRET", "", "Twitter app secret")

var MongoDB *mgo.Database

func init() {
	iniflags.Parse()
	session, err := mgo.Dial(*mongoCS)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	MongoDB = session.DB(*mongoName)

}

func ConfigureSwagger(wsContainer *restful.Container) {
	swaggerConfig := swagger.Config{
		WebServices:     wsContainer.RegisteredWebServices(), // you control what services are visible
		WebServicesUrl:  "http://" + *HttpBinding,
		ApiPath:         "/apidocs.json",
		SwaggerPath:     "/apidocs/",
		SwaggerFilePath: "./static/swagger",
	}
	swagger.RegisterSwaggerService(swaggerConfig, wsContainer)
}
