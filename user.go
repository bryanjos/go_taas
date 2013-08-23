package main

import (
	"errors"
	"github.com/emicklei/go-restful"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
)

type User struct {
	Id    string `json:"id,omitempty"`
	Email string
}

func (u User) validate() error {
	if u.Email == "" {
		return errors.New("Invalid Email")
	}

	return nil
}

type UserService struct {
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
		Doc("update a user").
		Param(ws.HeaderParameter("X-API-KEY", "api key").DataType("string")).
		Param(ws.BodyParameter("User", "representation of a user").DataType("main.User")).
		Reads(User{}))

	ws.Route(ws.PUT("").To(u.create).
		Doc("create a user").
		Param(ws.BodyParameter("User", "representation of a user").DataType("main.User")).
		Reads(User{}))

	ws.Route(ws.DELETE("").To(u.remove).
		Doc("delete a user").
		Param(ws.HeaderParameter("X-API-KEY", "api key").DataType("string")))

	restful.Add(ws)
}

func (u UserService) find(request *restful.Request, response *restful.Response) {
	session, err := mgo.Dial(DBSERVERNAME)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
	} else {
		defer session.Close()
		c := session.DB(DBNAME).C(USERCOLLECTION)

		api_key := request.HeaderParameter(APIKEYHEADER)

		user := User{}

		err = c.Find(bson.M{"id": api_key}).One(&user)
		if err != nil {
			response.WriteError(http.StatusNotFound, err)
		} else {
			response.WriteEntity(user)
		}
	}
}

func (u UserService) update(request *restful.Request, response *restful.Response) {
	session, err := mgo.Dial(DBSERVERNAME)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
	} else {
		defer session.Close()
		c := session.DB(DBNAME).C(USERCOLLECTION)
		api_key := request.HeaderParameter(APIKEYHEADER)

		usr := new(User)

		if err := request.ReadEntity(&usr); err == nil {
			user := User{}

			if err = c.Find(bson.M{"id": api_key}).One(&user); err != nil {
				response.WriteErrorString(http.StatusNotFound, "User Not Found")
			} else {
				user.Email = usr.Email

				if err = user.validate(); err != nil {
					response.WriteError(http.StatusPreconditionFailed, err)
				} else {
					err = c.Update(bson.M{"id": api_key}, &user)
					if err != nil {
						response.WriteError(http.StatusInternalServerError, err)
					} else {
						response.WriteErrorString(http.StatusOK, "OK")
					}
				}

			}
		} else {
			response.WriteError(http.StatusInternalServerError, err)
		}
	}
}

func (u UserService) create(request *restful.Request, response *restful.Response) {
	session, err := mgo.Dial(DBSERVERNAME)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
	} else {
		defer session.Close()
		c := session.DB(DBNAME).C(USERCOLLECTION)
		usr := new(User)

		if err := request.ReadEntity(&usr); err == nil {
			user := User{}

			if err = c.Find(bson.M{"Email": usr.Email}).One(&user); err != nil {

				if err = usr.validate(); err != nil {
					response.WriteError(http.StatusPreconditionFailed, err)
				} else {
					err = c.Insert(&user)
					if err != nil {
						response.WriteError(http.StatusInternalServerError, err)
					} else {
						_ = c.Find(bson.M{"Email": user.Email}).One(&user)
						response.WriteEntity(user)
					}
				}
			} else {
				response.WriteErrorString(http.StatusPreconditionFailed, "Email Already In Use")
			}
		} else {
			response.WriteError(http.StatusInternalServerError, err)
		}
	}
}

func (u UserService) remove(request *restful.Request, response *restful.Response) {
	session, err := mgo.Dial(DBSERVERNAME)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
	} else {
		defer session.Close()
		c := session.DB(DBNAME).C(USERCOLLECTION)
		api_key := request.HeaderParameter(APIKEYHEADER)
		user := User{}

		if err := c.Find(bson.M{"id": api_key}).One(&user); err != nil {
			response.WriteErrorString(http.StatusNotFound, "User Not Found")
		} else {
			err = c.Remove(bson.M{"id": api_key})
			if err != nil {
				response.WriteErrorString(http.StatusInternalServerError, "Error while deleting the user")
			} else {
				response.WriteErrorString(http.StatusOK, "")
			}

		}
	}
}
