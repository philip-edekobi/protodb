package server

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	gocask "github.com/philip-edekobi/go-cask"
)

type Server struct {
	Port       string
	DBInstance *gocask.BitCaskHandle
	DBIndex    *gocask.BitCaskHandle
}

func Newserver(port, dbLocation, indexLocation string) (*Server, error) {
	s := &Server{}
	s.Port = port

	db, err := gocask.Open(dbLocation)
	if err != nil {
		return nil, err
	}

	idx, err := gocask.Open(indexLocation)
	if err != nil {
		return nil, err
	}

	s.DBInstance = db
	s.DBIndex = idx

	return s, nil
}

func (s Server) Close() {
	s.DBIndex.Close()
	s.DBInstance.Close()
}

func (s Server) AddDoc(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	dec := json.NewDecoder(r.Body)
	var document map[string]any

	err := dec.Decode(&document)
	if err != nil {
		jsonResponse(w, nil, err)
		return
	}

	// new unique id for doc
	id := uuid.New().String()

	jsonStr, err := json.Marshal(document)
	if err != nil {
		jsonResponse(w, nil, err)
		return
	}

	err = s.DBInstance.Set(id, string(jsonStr))
	if err != nil {
		jsonResponse(w, nil, err)
		return
	}

	jsonResponse(w, map[string]any{
		"id": id,
	}, nil)
}

func (s Server) GetDoc(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}

func (s Server) SearchDocs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}
