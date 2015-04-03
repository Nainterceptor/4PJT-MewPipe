package configs

import (
    "github.com/vharitonsky/iniflags"
    "flag"
)

func Parse() {
    iniflags.Parse()
}

var staticPath = flag.String("static_path", "static", "Localisation for static files")

var HttpBinding = flag.String("http_binding", ":8080", "IP/Port to listen HTTP Server")

var MongoCS = flag.String("mongodb_CS", "localhost", "Connection endpoint for mongodb driver")
var MongoDB = flag.String("mongodb_DB", "MewPipe", "Database to mount")