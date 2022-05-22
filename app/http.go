package app

import (
	"github.com/gorilla/mux"
	"net/http"

	"skud/service"
)

// NewHTTPRouter creates new gorilla/mux router and registers all needed routes.
func NewHTTPRouter(svc *service.SkudService) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/access", HandleCheckAccess(svc)).Methods(http.MethodPost)
	return r
}
