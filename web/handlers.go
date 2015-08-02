package web

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func ApiV1Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "API v1 Index")
}

func ApiV1CreateJob(w http.ResponseWriter, r *http.Request) {
	//TODO
}
func ApiV1GetJob(w http.ResponseWriter, r *http.Request) {
	//TODO
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Fprintln(w, "Job show:", id)
}
func ApiV1DeleteJob(w http.ResponseWriter, r *http.Request) {
	//TODO
}

/*
func TodoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoId := vars["todoId"]
	fmt.Fprintln(w, "Todo show:", todoId)
}
*/
