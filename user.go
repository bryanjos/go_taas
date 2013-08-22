type User struct {
	Id     string
	Secret string
	Email  string
}

type UserService struct {
	session *r.Session
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
		Param(ws.HeaderParameter("X-API-SECRET", "api secret").DataType("string")).
		Writes(User{})) // on the response

	ws.Route(ws.POST("").To(u.update).
		// docs
		Doc("update a user").
		Param(ws.HeaderParameter("X-API-KEY", "api key").DataType("string")).
		Param(ws.HeaderParameter("X-API-SECRET", "api secret").DataType("string")).
		Param(ws.BodyParameter("User", "representation of a user").DataType("main.User")).
		Reads(User{})) // from the request

	ws.Route(ws.PUT("").To(u.create).
		// docs
		Doc("create a user").
		Param(ws.HeaderParameter("X-API-KEY", "api key").DataType("string")).
		Param(ws.HeaderParameter("X-API-SECRET", "api secret").DataType("string")).
		Param(ws.BodyParameter("User", "representation of a user").DataType("main.User")).
		Reads(User{})) // from the request

	ws.Route(ws.DELETE("").To(u.remove).
		// docs
		Doc("delete a user").
		Param(ws.HeaderParameter("X-API-KEY", "api key").DataType("string")).
		Param(ws.HeaderParameter("X-API-SECRET", "api secret").DataType("string")))

	restful.Add(ws)
}

func (u UserService) find(request *restful.Request, response *restful.Response) {
	api_key := request.HeaderParameter("X-API-KEY")
	api_secret := request.HeaderParameter("X-API-SECRET")

	var user User
	err := r.Table("user").Filter(r.Row.Attr("Id").Eq(api_key)).Filter(r.Row.Attr("Secret").Eq(api_secret)).Run(u.session).One(&user)

	if err != nil {
		response.WriteErrorString(http.StatusNotFound, "User Not Found")
	} else {
		response.WriteEntity(user)
	}
}

func (u UserService) update(request *restful.Request, response *restful.Response) {
	api_key := request.HeaderParameter("X-API-KEY")
	api_secret := request.HeaderParameter("X-API-SECRET")

	usr := new(User)
	err := request.ReadEntity(&usr)
	if err == nil {
		var user User
		err = r.Table("user").Filter(r.Row.Attr("Id").Eq(api_key)).Filter(r.Row.Attr("Secret").Eq(api_secret)).Run(u.session).One(&user)

		if err != nil {
			response.WriteErrorString(http.StatusNotFound, "User Not Found")
		} else {
			user.Email = usr.Email

			err = r.Table("user").Update(user).Run(u.session).One(&user)

			if err != nil {
				response.WriteError(http.StatusInternalServerError, err)
			} else {
				response.WriteEntity(user)
			}

		}
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}

func (u UserService) create(request *restful.Request, response *restful.Response) {
}

func (u UserService) remove(request *restful.Request, response *restful.Response) {
}