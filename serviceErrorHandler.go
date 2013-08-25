package main

import (
	"errors"
	"github.com/emicklei/go-restful"
	"github.com/op/go-logging"
	"net/http"
)

type ServiceErrorHandler struct {
	log *logging.Logger
}

func (s ServiceErrorHandler) WriteError(response *restful.Response, err error, status int, responseString string) {
	s.log.Error(err.Error())
	response.WriteServiceError(status, restful.ServiceError{status, responseString})
}

func (s ServiceErrorHandler) WriteInvalidAPIKeyError(response *restful.Response) {
	err := errors.New("Invalid API Key")
	s.WriteError(response, err, http.StatusPreconditionFailed, err.Error())
}

func (s ServiceErrorHandler) WriteInvalidPetKeyError(response *restful.Response) {
	err := errors.New("Invalid Pet Key")
	s.WriteError(response, err, http.StatusPreconditionFailed, err.Error())
}

func (s ServiceErrorHandler) WriteDatabaseConnectionError(response *restful.Response, err error) {
	s.WriteError(response, err, http.StatusInternalServerError, "Error Connecting to Database")
}

func (s ServiceErrorHandler) WriteNotFoundError(response *restful.Response, err error) {
	s.WriteError(response, err, http.StatusNotFound, "Not Found")
}

func (s ServiceErrorHandler) WriteInputError(response *restful.Response, err error) {
	s.WriteError(response, err, http.StatusInternalServerError, "Error Reading Input")
}

func (s ServiceErrorHandler) WriteUniqueEmailError(response *restful.Response) {
	err := errors.New("Email Already In Use")
	s.WriteError(response, err, http.StatusPreconditionFailed, err.Error())
}

func (s ServiceErrorHandler) WriteInvalidEmailError(response *restful.Response, err error) {
	s.WriteError(response, err, http.StatusPreconditionFailed, err.Error())
}

func (s ServiceErrorHandler) WriteInvalidEggError(response *restful.Response, err error) {
	s.WriteError(response, err, http.StatusPreconditionFailed, err.Error())
}

func (s ServiceErrorHandler) WriteUpdateError(response *restful.Response, err error) {
	s.WriteError(response, err, http.StatusInternalServerError, "Error During Update")
}

func (s ServiceErrorHandler) WriteCreateError(response *restful.Response, err error) {
	s.WriteError(response, err, http.StatusInternalServerError, "Error During Create")
}

func (s ServiceErrorHandler) WriteDeleteError(response *restful.Response, err error) {
	s.WriteError(response, err, http.StatusInternalServerError, "Error During Delete")
}
