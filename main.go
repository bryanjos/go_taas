package main

import (
	r "github.com/christopherhesse/rethinkgo"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"log"
	"net/http"
)

var sessionArray []*r.Session

func initDb() {
	session, err := r.Connect("localhost:28015", "taas")
	if err != nil {
		log.Fatal(err)
		return
	}

	err = r.DbCreate("taas").Run(session).Exec()
	if err != nil {
		log.Println(err)
	}

	err = r.TableCreate("user").Run(session).Exec()
	if err != nil {
		log.Println(err)
	}

	err = r.TableCreate("pet").Run(session).Exec()
	if err != nil {
		log.Println(err)
	}

	sessionArray = append(sessionArray, session)
}

func main() {

	initDb()

	u := UserService{DB{sessionArray[len(sessionArray)-1], "user"}}
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
