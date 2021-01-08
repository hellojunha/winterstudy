package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

var wd string


func init() {
	f, err := os.OpenFile("/home/joona0825/winterstudy.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Log file not found!")
	} else {
		log.SetOutput(f)
	}

	log.Println("instance is now running!")

	_wd, _ := os.Getwd()
	wd = _wd
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/{page:[0-9]+}", home)
	r.HandleFunc("/code", _registerCode).Methods(http.MethodPost)
	r.HandleFunc("/code/{id:[0-9]+}", _getCode).Methods(http.MethodGet)
	r.HandleFunc("/comment", _registerComment).Methods(http.MethodPost)
	r.HandleFunc("/category", _listCategory).Methods(http.MethodGet) // Listing All Categories
	r.HandleFunc("/category/{cat:[a-zA-Z0-9_]+", _getCategoryCode).Methods(http.MethodGet) // Get Categories' Code
	err := http.ListenAndServe(":9927", r)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func die(w http.ResponseWriter) {
	http.Error(w, "bad request", http.StatusBadRequest)
}

func home(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page, err := strconv.Atoi(vars["page"])

	var list []Post
	if err == nil {
		list = getPostList(page)
	} else {
		list = getPostList(0)
	}

	t, err := template.ParseFiles(wd + "/html/index.html")
	if err == nil {
		data := struct {
			Posts []Post
		} {
			Posts: list,
		}
		err := t.Execute(w, data)
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		log.Println(err.Error())
	}
}

func _registerCode(w http.ResponseWriter, r *http.Request) {
	registerPost(r.PostFormValue("captcha"), r.PostFormValue("code"), r.PostFormValue("category"))
}

func _getCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err == nil {
		post := getPost(id)
		t, err := template.ParseFiles(wd + "/html/post.html")
		if err == nil {
			err := t.Execute(w, post)
			if err != nil {
				log.Println(err.Error())
			}
		} else {
			log.Println(err.Error())
		}
	} else {
		die(w)
		return
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

// Listing All Categories
func _ListCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if err == nil {
		categroies := _getCategoryList() // returns string type slice 
		t, err := template.ParseFiles(wd + "/html/categoryList.html")

		if err == nil {
			err := t.Execute(w, categories) // check this code and categoryList.html
			if err != nil {
				log.Println(err.Error())
			}
		} else {
			log.Println(err.Error())
		}
	} else {
		die(w)
		return
	}

}

// Render codes matches category in DB
func _getCategoryCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["cat"]

	list = getPostListFromCategory(cat)

	t, err := template.ParseFiles(wd + "/html/index.html")
	if err == nil {
		data := struct {
			Posts []Post
		} {
			Posts: list,
		}
		err := t.Execute(w, data)
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		log.Println(err.Error())
	}
}