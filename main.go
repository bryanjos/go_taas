package main

import (
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"github.com/op/go-logging"
	"net/http"
)

//TODO: Move to Mongo

const APIKEYHEADER = "X-API-KEY"
const PETKEYHEADER = "X-PET-KEY"
const DBSERVERNAME = "localhost"
const DBNAME = "taas"
const USERCOLLECTION = "user"
const PETCOLLECTION = "pet"

func main() {
	var log = logging.MustGetLogger(DBNAME)
	logging.SetLevel(logging.INFO, DBNAME)

	serviceErrorHandler := new(ServiceErrorHandler)
	serviceErrorHandler.log = log

	u := UserService{serviceErrorHandler}
	u.Register()

	p := PetService{serviceErrorHandler}
	p.Register()

	config := swagger.Config{
		WebServices:    restful.RegisteredWebServices(), // you control what services are visible
		WebServicesUrl: "http://localhost:5000",
		ApiPath:        "/apidocs.json",

		// Optionally, specifiy where the UI is located
		SwaggerPath:     "/",
		SwaggerFilePath: "/Users/bryanjos/Projects/go_taas/swagger-ui",
	}

	swagger.InstallSwaggerService(config)

	log.Info("start listening on localhost:5000")
	log.Fatal(http.ListenAndServe(":5000", nil))

}
