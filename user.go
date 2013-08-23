package main

import (
	"errors"
	"github.com/emicklei/go-restful"
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
	db DB
}

func (u UserService) Register() {
	ws := new(restful.WebService)
	ws.
		Path("/user").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML) // you can specify this per route as well

	ws.Route(ws.GET("").To(u.find).
		// docs
		Doc("get a user").
		Param(ws.HeaderParameter("X-API-KEY", "api key").DataType("string")).
		Writes(User{})) // on the response

	ws.Route(ws.POST("").To(u.update).
		// docs
		Doc("update a user").
		Param(ws.HeaderParameter("X-API-KEY", "api key").DataType("string")).
		Param(ws.BodyParameter("User", "representation of a user").DataType("main.User")).
		Reads(User{})) // from the request

	ws.Route(ws.PUT("").To(u.create).
		// docs
		Doc("create a user").
		Param(ws.BodyParameter("User", "representation of a user").DataType("main.User")).
		Reads(User{})) // from the request

	ws.Route(ws.DELETE("").To(u.remove).
		// docs
		Doc("delete a user").
		Param(ws.HeaderParameter("X-API-KEY", "api key").DataType("string")))

	restful.Add(ws)
}

func (u UserService) find(request *restful.Request, response *restful.Response) {
	api_key := request.HeaderParameter("X-API-KEY")

	var user User
	err := u.db.find("id", api_key, &user)

	if err != nil {
		response.WriteError(http.StatusNotFound, err)
	} else {
		response.WriteEntity(user)
	}
}

func (u UserService) update(request *restful.Request, response *restful.Response) {
	api_key := request.HeaderParameter("X-API-KEY")

	usr := new(User)
	err := request.ReadEntity(&usr)
	if err == nil {
		var user User
		err = u.db.find("id", api_key, &user)

		if err != nil {
			response.WriteErrorString(http.StatusNotFound, "User Not Found")
		} else {
			user.Email = usr.Email
			err = user.validate()

			if err != nil {
				response.WriteError(http.StatusPreconditionFailed, err)
			} else {
				u.db.update(user)
				response.WriteErrorString(http.StatusOK, "")
			}

		}
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}

func (u UserService) create(request *restful.Request, response *restful.Response) {
	usr := new(User)
	err := request.ReadEntity(&usr)
	if err == nil {
		var user User
		err = u.db.find("Email", usr.Email, &user)

		if err != nil {
			err = usr.validate()

			if err != nil {
				response.WriteError(http.StatusPreconditionFailed, err)
			} else {
				u.db.create(usr)
				err = u.db.find("Email", usr.Email, &user)
				if err != nil {
					response.WriteError(http.StatusInternalServerError, err)
				} else {
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

func (u UserService) remove(request *restful.Request, response *restful.Response) {
	api_key := request.HeaderParameter("X-API-KEY")

	var user User
	err := u.db.find("id", api_key, &user)

	if err != nil {
		response.WriteErrorString(http.StatusNotFound, "User Not Found")
	} else {
		u.db.delete("id", api_key)
		response.WriteErrorString(http.StatusOK, "")
	}
}
