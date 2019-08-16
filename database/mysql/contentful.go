package mysql

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/gorilla/mux"
	"github.com/jivalabs/picsum-photos/api/handler"
	"github.com/jivalabs/picsum-photos/database"
	"github.com/jivalabs/picsum-photos/logger"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"time"
)

const HeaderContentType = "Content-Type"

//const MIMEApplicationJSON = "application/json"
const MIMEApplicationJSON = "application/json; charset=utf-8"

// Content Management API
//const apiBaseUrl = "https://api.contentful.com"
const managementApiBaseUrl = "http://localhost:8080"

// Access token name
const tokenName = "Example Key 1"

// Access token description
const tokenDescription = "We’ve created an example API key for you to help you get started."

// Space ID
const spaceID = "6yry1jvchfmm"

// Content Delivery API access tokenl
const cdApiToken = "5QTPUOdMS-f6eI99sWI2a4MhLf-CP9VRhFEqrxcjks4"

// Content Preview API access token
const cpApiToken = "rCu3Ta1s2d1HMlQHO6EwvNH6xXTwHRwsBceJD9yUTTw"

// Environment
const environment = "master"

// contentful api related structs

type SpaceApi struct {
	Database database.SpaceProvider
	Log      *logger.Logger
}

func logFields(r *http.Request, keysAndValues ...interface{}) []interface{} {
	ctx := r.Context()
	id := handler.GetReqID(ctx)

	return append([]interface{}{"request-id", id}, keysAndValues...)
}

func (s *SpaceApi) logError(r *http.Request, message string, err error) {
	s.Log.Errorw(message, logFields(r, "error", err)...)
}

type link struct {
	a string
}

type sys struct {
	Type_       string    `json:"type"`
	Id          string    `json:"id"`
	Space       link      `json:"space"`
	ContentType string    `json:"contentType"`
	Revision    int       `json:"revision"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Locale      string    `json:"locale"`
}

type response struct {
	Sys   sys           `json:"sys"`
	Skip  int           `json:"skip"`
	Limit int           `json:"limit"`
	Total int           `json:"total"`
	Items []interface{} `json:"items"`
}

type Symbol string
type Text string
type Details struct {
	size int
}

type file struct {
	FileName    Symbol  `json:"fileName"`
	ContentType Symbol  `json:"contentType"`
	Url         Symbol  `json:"url"`
	Details     Details `json:"details"`
}

type fields struct {
	title       string
	description string
	file        file
}

type asset struct {
	sys    sys
	fields fields
}

// Authentication

// get access token

func (s *SpaceApi) ContentfulRouter() http.Handler {
	return nil
}

type createSpaceT struct {
	Name          string `json:"Name"`
	DefaultLocale string `json:"defaultLocale"`
}

// intParam tries to get a param and convert it to an Integer
func intParam(r *http.Request, name string) (int, bool) {
	vars := mux.Vars(r)

	if val, ok := vars[name]; ok {
		val, err := strconv.Atoi(val)
		return val, err == nil
	}

	return -1, false
}

func (s *SpaceApi) Router(h http.Handler) http.Handler {
	r := chi.NewRouter()
	r.Mount("/spaces", s.spaceRouter())
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
	return r
}

func (s *SpaceApi) spaceRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/{spaceId}", s.getSpaceHandler)
	r.Get("/", s.getSpaceListHandler)
	r.Post("/", s.createSpaceHandler)

	return r
}

// ListAll returns a list of all the images
func (s *SpaceProvider) ListAll() ([]database.Space, error) {
	i := []database.Space{}
	err := s.db.Select(&i, "select * from space")

	if err != nil {
		return nil, err
	}

	return i, nil
}

func (s *SpaceApi) getSpaceListHandler(w http.ResponseWriter, r *http.Request) {
	spaces, err := s.Database.GetSpaceList()
	if err != nil {
		s.logError(r, "error getting list of spaces from database", err)
		return
	}
	sendJsonOk(w, spaces)
}

func (s *SpaceApi) getSpaceHandler(w http.ResponseWriter, r *http.Request) {
	sendJsonOk(w, "getSpaceHandler spaceById output")
}

// https://www.contentful.com/developers/docs/references/content-management-api/#/reference/spaces/spaces-collection/create-a-space/console/curl
func (s *SpaceApi) createSpaceHandler(w http.ResponseWriter, r *http.Request) {
	// cma_token from header
	// method POST
	// Content-Type: "application/json"
	// payload json createSpaceT

	var space createSpaceT

	err := json.NewDecoder(r.Body).Decode(&space)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = sendJsonOk(w, space)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func sendJsonOk(w http.ResponseWriter, obj interface{}) error {
	return sendJSON(w, http.StatusOK, obj)
}

func sendJSON(w http.ResponseWriter, status int, obj interface{}) error {
	w.Header().Set(HeaderContentType, MIMEApplicationJSON)
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error encoding json response: %v", obj))
	}
	w.WriteHeader(status)
	_, err = w.Write(jsonData)
	return err
}

// Content Delivery API

// Content Management API
