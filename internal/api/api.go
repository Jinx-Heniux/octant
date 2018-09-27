package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/heptio/developer-dash/internal/overview"
)

type errorResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type notFoundResponse struct {
	Error errorResponse `json:"error,omitempty"`
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	r := &notFoundResponse{
		Error: errorResponse{
			Code:    code,
			Message: message,
		},
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(r)
}

// API is the API for the dashboard client
type API struct {
	mux *mux.Router
}

var _ http.Handler = (*API)(nil)

// New creates an instance of API.
func New(prefix string, o overview.Interface) *API {
	router := mux.NewRouter()
	s := router.PathPrefix(prefix).Subrouter()

	namespacesService := newNamespaces(o)
	s.Handle("/namespaces", namespacesService)

	navigationService := &navigation{}
	s.Handle("/navigation", navigationService)

	contentService := &content{}
	s.Handle("/content/{path:.*}", contentService)

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)

		resp := &notFoundResponse{
			Error: errorResponse{
				Code:    http.StatusNotFound,
				Message: "not found",
			},
		}

		json.NewEncoder(w).Encode(resp)
	})

	return &API{
		mux: router,
	}
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}