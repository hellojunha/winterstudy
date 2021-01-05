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
	r.HandleFunc("/post", _registerCodeForm).Methods(http.MethodGet)
	r.HandleFunc("/post/register", _registerCode).Methods(http.MethodPost)
	r.HandleFunc("/code/{id:[0-9]+}", _getCode).Methods(http.MethodGet)
	r.HandleFunc("/comment", _registerComment).Methods(http.MethodPost)
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
	var pages int

	if err == nil {
		list, pages = getPostList(page)
	} else {
		list, pages = getPostList(0)
	}

	t, err := template.ParseFiles(wd + "/html/index.html")
	if err == nil {
		data := struct {
			Posts []Post
			Pages []int
		} {
			Posts: list,
			Pages: make([]int, pages),
		}
		err := t.Execute(w, data)
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		log.Println(err.Error())
	}
}

func _registerCodeForm(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(wd + "/html/post.html")
	if err == nil {
		data := struct {
			CaptchaKey string
		} {
			CaptchaKey: CAPTCHA_KEY,
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
	result := registerPost(r.PostFormValue("captcha"), r.PostFormValue("code"))
	if result == -1 {
		t, _ := template.ParseFiles(wd + "/html/post_fail.html")
		t.Execute(w, r.PostFormValue("code"))
		return
	}

	t, _ := template.ParseFiles(wd + "/html/post_success.html")
	t.Execute(w, result)
}

func _getCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err == nil {
		post := getPost(id)
		t, err := template.ParseFiles(wd + "/html/code.html")
		if err == nil {
			data := struct {
				Post Post
				CaptchaKey string
			} {
				Post: *post,
				CaptchaKey: CAPTCHA_KEY,
			}
			err := t.Execute(w, data)
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
	comment := r.PostFormValue("comment")
	if err == nil {
		result := registerComment(r.PostFormValue("captcha"), postId, comment)
		if result {
			http.Redirect(w, r, "https://study.alfr.kr/code/" + r.PostFormValue("post_id"), http.StatusSeeOther)
		} else {
			fmt.Fprintf(w, "Failed to register comment.\nThis was your comment:\n\n%s", comment)
		}
	} else {
		die(w)
	}

}
