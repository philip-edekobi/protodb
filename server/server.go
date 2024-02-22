package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

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

	s.index(id, document)

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

	isRange := false
	idsArgumentCount := map[string]int{}
	nonRangeArgs := 0

	for _, argument := range q.Ands {
		if argument.Op == "=" {
			nonRangeArgs++

			ids, err := s.lookup(
				fmt.Sprintf("%s=%v", strings.Join(argument.Key, "."), argument.Value),
			)
			if err != nil {
				jsonResponse(w, nil, err)
				return
			}

			for _, id := range ids {
				idsArgumentCount[id]++
			}
		} else {
			isRange = true
		}
	}

	var matchingIds []string
	for id, count := range idsArgumentCount {
		if count == nonRangeArgs {
			matchingIds = append(matchingIds, id)
		}
	}

	var documents []map[string]any
	if r.URL.Query().Get("skip-index") == "true" {
		matchingIds = nil
	}

	if len(matchingIds) > 0 {
		for _, id := range matchingIds {
			doc, err := s.docByID(id)
			if err != nil {
				jsonResponse(w, nil, err)
				return
			}

			if !isRange || q.Match(doc) {
				documents = append(documents, map[string]any{
					"id":   id,
					"body": doc,
				})
			}
		}
	} else {
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
	}

	jsonResponse(w, map[string]any{"documents": documents, "count": len(documents)}, nil)
}

func (s Server) index(id string, document map[string]any) {
	pv := getPathValues(document, "")

	for _, pathValue := range pv {
		idsString, err := s.DBIndex.Get(pathValue)
		if err != nil && err != gocask.ErrKeyNotFound {
			log.Printf("failed to look up path value [%#v]: %s", document, err)
		}

		if len(idsString) == 0 {
			idsString = id
		} else {
			ids := strings.Split(idsString, ",")

			found := false
			for _, existindId := range ids {
				if id == existindId {
					found = true
					break
				}
			}

			if !found {
				idsString += "," + id
			}
		}

		err = s.DBIndex.Set(pathValue, idsString)
		if err != nil {
			log.Printf("could not update index: %s", err)
		}
	}
}

func getPathValues(obj map[string]any, prefix string) []string {
	var pvs []string

	for key, val := range obj {
		switch t := val.(type) {
		case map[string]any:
			pvs = append(pvs, getPathValues(t, key)...)
		case []interface{}:
			// skip arrays
			continue
		}

		if prefix != "" {
			key = prefix + "." + key
		}

		pvs = append(pvs, fmt.Sprintf("%s=%v", key, val))
	}

	return pvs
}

func (s Server) lookup(pathValue string) ([]string, error) {
	idsString, err := s.DBIndex.Get(pathValue)
	if err != nil && err != gocask.ErrKeyNotFound {
		return nil, fmt.Errorf("could not lookup pathvalue [%#v]: %s", pathValue, err)
	}

	if len(idsString) == 0 {
		return nil, nil
	}

	return strings.Split(idsString, ","), nil
}
