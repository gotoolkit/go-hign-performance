package handler

import (
	"github.com/gorilla/mux"
	"github.com/llitfkitfk/GoHighPerformance/pkg/db"
	"net/http"
)

type CreateHandler struct {
	db db.DB
}

func NewCreateHandler(db db.DB) *CreateHandler {
	return &CreateHandler{db: db}
}

func (c *CreateHandler) RegisterRoute(r *mux.Router) {
	r.Handle("/test", c).Methods("POST")
}

func (c *CreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
