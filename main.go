package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/philip-edekobi/protodb/server"
)

const (
	DATADIR  = "./data/"
	INDEXDIR = "./data/index/"
	TESTDIR  = "./testData/"
	PORT     = "5477"
)

func main() {
	server := server.Newserver(PORT, DATADIR, INDEXDIR)

	router := httprouter.New()

	router.POST("/docs", server.AddDoc)
	router.GET("/docs", server.SearchDocs)
	router.GET("/docs/:id", server.GetDoc)

	log.Println("server is active on port:", server.Port)
	log.Fatal(http.ListenAndServe(":"+server.Port, router))
}
