package main

import (
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"log"
	"net/http"
)

//TODO: Move to Mongo

const APIKEYHEADER = "X-API-KEY"
const DBSERVERNAME = "localhost"
const DBNAME = "taas"
const USERCOLLECTION = "user"

func main() {

	u := UserService{}
	u.Register()

	config := swagger.Config{
		WebServices:    restful.RegisteredWebServices(), // you control what services are visible
		WebServicesUrl: "http://localhost:5000",
		ApiPath:        "/apidocs.json",

		// Optionally, specifiy where the UI is located
		SwaggerPath:     "/apidocs/",
		SwaggerFilePath: "/Users/bryanjos/Projects/go_taas/swagger-ui",
	}

	swagger.InstallSwaggerService(config)

	log.Printf("start listening on localhost:5000")
	log.Fatal(http.ListenAndServe(":5000", nil))

}
