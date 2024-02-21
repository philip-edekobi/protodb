package server

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	gocask "github.com/philip-edekobi/go-cask"

	"github.com/philip-edekobi/protodb/searcher"
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

func (s Server) docByID(id string) (map[string]any, error) {
	val, err := s.DBInstance.Get(id)
	if err != nil {
		return nil, err
	}

	var document map[string]any

	err = json.Unmarshal([]byte(val), &document)

	return document, nil
}

func (s Server) GetDoc(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	doc, err := s.docByID(id)
	if err != nil {
		jsonResponse(w, nil, err)
		return
	}

	jsonResponse(w, map[string]any{
		"document": doc,
	}, nil)
}

func (s Server) SearchDocs(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	q, err := searcher.ParseFilter(r.URL.Query().Get("q"))
	if err != nil {
		jsonResponse(w, nil, err)
		return
	}

	var documents []map[string]any

	keys := s.DBInstance.ListKeys()

	for _, key := range keys {
		var doc map[string]any

		val, err := s.DBInstance.Get(key)
		if err != nil {
			jsonResponse(w, nil, err)
			return
		}

		err = json.Unmarshal([]byte(val), &doc)
		if err != nil {
			jsonResponse(w, nil, err)
			return
		}

		if q.Match(doc) {
			documents = append(documents, map[string]any{
				"id":   key,
				"body": doc,
			})
		}
	}

	jsonResponse(w, map[string]any{"documents": documents, "count": len(documents)}, nil)
}
