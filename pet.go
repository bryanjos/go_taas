package main

import (
	"crypto/sha512"
	"errors"
	"github.com/emicklei/go-restful"
	"io"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/url"
	"time"
)

type Pet struct {
	Id           bson.ObjectId `json:"Id"           bson:"_id"`
	UserId       bson.ObjectId `json:"UserId"       bson:"userId"`
	Name         string
	Url          string
	Health       int
	Stamina      int
	Attitude     int
	Intelligence int
	Age          int
	Gender       string
	Created      time.Time
	Updated      time.Time
}

type Egg struct {
	Name string
	Url  string
}

func (e Egg) validate() error {
	if _, err := url.Parse(e.Url); e.Url == "" || err != nil {
		return errors.New("Invalid Url")
	}

	if e.Name == "" {
		return errors.New("Invalid Name")
	}

	return nil
}

func (e Egg) hatch(api_key bson.ObjectId) Pet {
	hash := sha512.New()
	io.WriteString(hash, e.Url)
	byteArray := hash.Sum(nil)

	pet := Pet{}
	pet.Id = bson.NewObjectId()
	pet.UserId = api_key
	pet.Url = e.Url
	pet.Name = e.Name
	pet.Health = int(byteArray[0]) % 100
	pet.Stamina = int(byteArray[1]) % 100
	pet.Attitude = int(byteArray[2]) % 100
	pet.Intelligence = int(byteArray[3]) % 100
	pet.Age = int(byteArray[4]) % 100
	if int(byteArray[4]) > 100 {
		pet.Gender = "M"
	} else {
		pet.Gender = "F"
	}

	pet.Created = time.Now()
	pet.Updated = time.Now()

	return pet
}

type PetService struct {
	errorHandler *ServiceErrorHandler
}

func (p PetService) Register() {
	ws := new(restful.WebService)
	ws.
		Path("/pet").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.GET("").To(p.find).
		Doc("get your pet").
		Param(ws.HeaderParameter("X-API-KEY", "api key").DataType("string")).
		Param(ws.HeaderParameter("X-PET-KEY", "pet key").DataType("string")).
		Writes(Pet{}))

	ws.Route(ws.PUT("").To(p.create).
		Doc("create a pet").
		Param(ws.HeaderParameter("X-API-KEY", "api key").DataType("string")).
		Param(ws.BodyParameter("Egg", "Json with url and name attributes").DataType("Egg")).
		Writes(Pet{}))

	ws.Route(ws.DELETE("").To(p.remove).
		Doc("delete a pet").
		Param(ws.HeaderParameter("X-API-KEY", "api key").DataType("string")).
		Param(ws.HeaderParameter("X-PET-KEY", "pet key").DataType("string")).
		Writes(""))

	restful.Add(ws)
}

func (p PetService) find(request *restful.Request, response *restful.Response) {
	if bson.IsObjectIdHex(request.HeaderParameter(APIKEYHEADER)) == false {
		p.errorHandler.WriteInvalidAPIKeyError(response)
		return
	}

	if bson.IsObjectIdHex(request.HeaderParameter(PETKEYHEADER)) == false {
		p.errorHandler.WriteInvalidPetKeyError(response)
		return
	}

	api_key := bson.ObjectIdHex(request.HeaderParameter(APIKEYHEADER))
	pet_key := bson.ObjectIdHex(request.HeaderParameter(PETKEYHEADER))
	session, err := mgo.Dial(DBSERVERNAME)

	if err != nil {
		p.errorHandler.WriteDatabaseConnectionError(response, err)
		return
	}

	defer session.Close()
	c := session.DB(DBNAME).C(PETCOLLECTION)

	pet := Pet{}

	if err = c.Find(bson.M{"_id": pet_key, "userId": api_key}).One(&pet); err != nil {
		p.errorHandler.WriteNotFoundError(response, err)
		return
	}

	response.WriteEntity(pet)
}

func (p PetService) create(request *restful.Request, response *restful.Response) {
	if bson.IsObjectIdHex(request.HeaderParameter(APIKEYHEADER)) == false {
		p.errorHandler.WriteInvalidAPIKeyError(response)
		return
	}

	api_key := bson.ObjectIdHex(request.HeaderParameter(APIKEYHEADER))
	session, err := mgo.Dial(DBSERVERNAME)

	if err != nil {
		p.errorHandler.WriteDatabaseConnectionError(response, err)
		return
	}

	defer session.Close()
	c := session.DB(DBNAME).C(PETCOLLECTION)
	egg := new(Egg)

	if err := request.ReadEntity(&egg); err != nil {
		p.errorHandler.WriteInputError(response, err)
		return
	}

	if err = egg.validate(); err != nil {
		p.errorHandler.WriteInvalidEggError(response, err)
		return
	}

	pet := egg.hatch(api_key)

	if err = c.Insert(&pet); err != nil {
		p.errorHandler.WriteCreateError(response, err)
		return
	}

	response.WriteEntity(pet)
}

func (p PetService) remove(request *restful.Request, response *restful.Response) {
	if bson.IsObjectIdHex(request.HeaderParameter(APIKEYHEADER)) == false {
		p.errorHandler.WriteInvalidAPIKeyError(response)
		return
	}

	if bson.IsObjectIdHex(request.HeaderParameter(PETKEYHEADER)) == false {
		p.errorHandler.WriteInvalidPetKeyError(response)
		return
	}

	api_key := bson.ObjectIdHex(request.HeaderParameter(APIKEYHEADER))
	pet_key := bson.ObjectIdHex(request.HeaderParameter(PETKEYHEADER))
	session, err := mgo.Dial(DBSERVERNAME)

	if err != nil {
		p.errorHandler.WriteDatabaseConnectionError(response, err)
		return
	}

	defer session.Close()
	c := session.DB(DBNAME).C(PETCOLLECTION)

	if err = c.Remove(bson.M{"_id": pet_key, "userId": api_key}); err != nil {
		p.errorHandler.WriteDeleteError(response, err)
		return
	}

	response.WriteEntity("OK")
}
