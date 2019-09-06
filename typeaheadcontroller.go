package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type TypeaheadController struct {
	Port    int
	service *TypeaheadService
}

func NewAccountController(port int, service *TypeaheadService) *TypeaheadController {
	return &TypeaheadController{
		Port:    port,
		service: service,
	}
}

func (controller TypeaheadController) Start() {
	router := mux.NewRouter()
	router.HandleFunc("/typeahead", controller.QueryHandler).Methods("GET")
	router.HandleFunc("/typeahead/", controller.QueryHandler).Methods("GET")
	router.HandleFunc("/typeahead/{prefix}", controller.QueryHandler).Methods("GET")
	router.HandleFunc("/typeahead/set", controller.VoteHandler).Methods("POST")

	loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	http.ListenAndServe(":"+strconv.Itoa(controller.Port), loggedRouter)
}

func (controller TypeaheadController) QueryHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	prefix, prefixFound := vars["prefix"]

	log.Println("Prefix:", prefix)

	results := controller.service.Query(prefix)

	if prefixFound {
		shotEntry, foundShot := controller.service.Get(prefix)

		if foundShot {
			entryToRemove := -1

			for i, v := range results {
				if v.Name == shotEntry.Name {
					entryToRemove = i
					log.Println("Entry to remove:", v.Name)
				}
			}

			if entryToRemove != -1 {
				results = append(results[:entryToRemove], results[entryToRemove+1:]...)
			}

			results = append([]Result{*shotEntry}, results...)
		}
	}

	js, err := json.Marshal(results)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	w.Write([]byte("\n"))
}

type VoteRequest struct {
	Name string `json:name`
}

func (controller TypeaheadController) VoteHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var voteRequest VoteRequest

	err := decoder.Decode(&voteRequest)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = controller.service.Vote(voteRequest.Name)

	if err != nil {
		http.Error(w, "Name not found", http.StatusBadRequest)
		return
	}

	result, _ := controller.service.Get(voteRequest.Name)

	js, err := json.Marshal(*result)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	w.Write([]byte("\n"))
}
