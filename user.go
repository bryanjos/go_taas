package main

import (
	"errors"
	"github.com/emicklei/go-restful"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"strings"
	"time"
)

type User struct {
	Id      bson.ObjectId `json:"Id"           bson:"_id"`
	Email   string
	Created time.Time
	Updated time.Time
}

type Email struct {
	Email string
}

func (u Email) validate() error {
	if u.Email == "" || strings.Contains(u.Email, "@") == false {
		return errors.New("Invalid Email")
	}

	return nil
}

type UserService struct {
	errorHandler *ServiceErrorHandler
}

func (u UserService) Register() {
	ws := new(restful.WebService)
	ws.
		Path("/user").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.GET("").To(u.find).
		Doc("get a user").
		Param(ws.HeaderParameter("X-API-KEY", "api key").DataType("string")).
		Writes(User{}))

	ws.Route(ws.POST("").To(u.update).
		Doc("updates a user's email").
		Param(ws.HeaderParameter("X-API-KEY", "api key").DataType("string")).
		Param(ws.BodyParameter("Email", "json with an email attribute").DataType("Email")).
		Writes(User{}))

	ws.Route(ws.PUT("").To(u.create).
		Doc("create a user").
		Param(ws.BodyParameter("Email", "json with an email attribute").DataType("Email")).
		Writes(User{}))

	ws.Route(ws.DELETE("").To(u.remove).
		Doc("delete a user").
		Param(ws.HeaderParameter("X-API-KEY", "api key").DataType("string")).
		Writes(""))

	restful.Add(ws)
}

func (u UserService) find(request *restful.Request, response *restful.Response) {
	if bson.IsObjectIdHex(request.HeaderParameter(APIKEYHEADER)) == false {
		u.errorHandler.WriteInvalidAPIKeyError(response)
		return
	}

	api_key := bson.ObjectIdHex(request.HeaderParameter(APIKEYHEADER))
	session, err := mgo.Dial(DBSERVERNAME)

	if err != nil {
		u.errorHandler.WriteDatabaseConnectionError(response, err)
		return
	}

	defer session.Close()
	c := session.DB(DBNAME).C(USERCOLLECTION)

	user := User{}

	if err = c.FindId(api_key).One(&user); err != nil {
		u.errorHandler.WriteNotFoundError(response, err)
		return
	}

	response.WriteEntity(user)
}

func (u UserService) update(request *restful.Request, response *restful.Response) {
	if bson.IsObjectIdHex(request.HeaderParameter(APIKEYHEADER)) == false {
		u.errorHandler.WriteInvalidAPIKeyError(response)
		return
	}

	api_key := bson.ObjectIdHex(request.HeaderParameter(APIKEYHEADER))
	session, err := mgo.Dial(DBSERVERNAME)

	if err != nil {
		u.errorHandler.WriteDatabaseConnectionError(response, err)
		return
	}

	defer session.Close()
	c := session.DB(DBNAME).C(USERCOLLECTION)

	usr := new(Email)

	if err := request.ReadEntity(&usr); err != nil {
		u.errorHandler.WriteInputError(response, err)
		return
	}

	if err = usr.validate(); err != nil {
		u.errorHandler.WriteInvalidEmailError(response, err)
		return
	}

	if count, _ := c.Find(bson.M{"email": usr.Email}).Count(); count > 0 {
		u.errorHandler.WriteUniqueEmailError(response)
		return
	}

	user := User{}
	user.Email = usr.Email
	user.Updated = time.Now()

	if c.FindId(api_key).One(&user); err != nil {
		u.errorHandler.WriteNotFoundError(response, err)
		return
	}

	if err = c.UpdateId(api_key, &user); err != nil {
		u.errorHandler.WriteUpdateError(response, err)
		return
	}

	response.WriteEntity(user)

}

func (u UserService) create(request *restful.Request, response *restful.Response) {
	session, err := mgo.Dial(DBSERVERNAME)
	if err != nil {
		u.errorHandler.WriteDatabaseConnectionError(response, err)
		return
	}

	defer session.Close()
	c := session.DB(DBNAME).C(USERCOLLECTION)
	usr := new(Email)

	if err := request.ReadEntity(&usr); err != nil {
		u.errorHandler.WriteInputError(response, err)
		return
	}

	if count, _ := c.Find(bson.M{"email": usr.Email}).Count(); count > 0 {
		u.errorHandler.WriteUniqueEmailError(response)
		return
	}

	if err = usr.validate(); err != nil {
		u.errorHandler.WriteInvalidEmailError(response, err)
		return
	}

	user := User{bson.NewObjectId(), usr.Email, time.Now(), time.Now()}

	if err = c.Insert(&user); err != nil {
		u.errorHandler.WriteCreateError(response, err)
		return
	}

	response.WriteEntity(user)
}

func (u UserService) remove(request *restful.Request, response *restful.Response) {
	if bson.IsObjectIdHex(request.HeaderParameter(APIKEYHEADER)) == false {
		u.errorHandler.WriteInvalidAPIKeyError(response)
		return
	}

	api_key := bson.ObjectIdHex(request.HeaderParameter(APIKEYHEADER))
	session, err := mgo.Dial(DBSERVERNAME)

	if err != nil {
		u.errorHandler.WriteDatabaseConnectionError(response, err)
		return
	}

	defer session.Close()
	c := session.DB(DBNAME).C(USERCOLLECTION)

	if err = c.RemoveId(api_key); err != nil {
		u.errorHandler.WriteDeleteError(response, err)
		return
	}

	response.WriteEntity("OK")
}
