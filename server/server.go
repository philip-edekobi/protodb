package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	gocask "github.com/philip-edekobi/go-cask"
)

type Server struct {
	Port       string
	DBInstance *gocask.BitCaskHandle
	DBIndex    *gocask.BitCaskHandle
}

func Newserver(port, dbLocation, indexLocation string) *Server {
	s := &Server{}
	s.Port = port

	db := gocask.Open(dbLocation)
	idx := gocask.Open(indexLocation)

	s.DBInstance = db
	s.DBIndex = idx

	return s
}

func (s Server) AddDoc(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
}

func (s Server) GetDoc(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}

func (s Server) SearchDocs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}
