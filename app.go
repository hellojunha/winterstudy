package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strconv"
)

func init() {
	f, err := os.OpenFile("/home/joona0825/winterstudy.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Log file not found!")
	} else {
		log.SetOutput(f)
	}

	log.Println("instance is now running!")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/{page:[0-9]+}", home)
	r.HandleFunc("/code", _registerCode).Methods(http.MethodPost)
	r.HandleFunc("/code/{id:[0-9]+}", _getCode).Methods(http.MethodGet)
	r.HandleFunc("/comment", _registerComment).Methods(http.MethodPost)
	http.ListenAndServe(":9927", r)
}

func die(w http.ResponseWriter) {
	http.Error(w, "bad request", http.StatusBadRequest)
}

func home(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page, err := strconv.Atoi(vars["page"])
	if err == nil {
		getPostList(page)
	} else {
		getPostList(0)
	}
}

func _registerCode(w http.ResponseWriter, r *http.Request) {
	registerPost(r.PostFormValue("captcha"), r.PostFormValue("code"))
}

func _getCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err == nil {
		getPost(id)
	} else {
		die(w)
	}
}

func _registerComment(w http.ResponseWriter, r *http.Request) {
	postId, err := strconv.Atoi(r.PostFormValue("post_id"))
	if err == nil {
		registerComment(r.PostFormValue("captcha"), postId, r.PostFormValue("text"))
	} else {
		die(w)
	}

}
