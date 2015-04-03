package configs

import (
    "github.com/vharitonsky/iniflags"
    "flag"
    "gopkg.in/mgo.v2"
)

func Parse() {
    iniflags.Parse()
}

var staticPath = flag.String("static_path", "static", "Localisation for static files")

var HttpBinding = flag.String("http_binding", ":8080", "IP/Port to listen HTTP Server")

var mongoCS = flag.String("mongodb_CS", "localhost", "Connection endpoint for mongodb driver")
var mongoName = flag.String("mongodb_DB", "MewPipe", "Database to mount")

var MongoDB = getMongoDBVar()


func getMongoDBVar() *mgo.Database {
    session, err := mgo.Dial(*mongoCS)
    if err != nil {
        panic(err)
    }
    session.SetMode(mgo.Monotonic, true)
    return session.DB(*mongoName)
}