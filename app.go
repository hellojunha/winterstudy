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
	r.HandleFunc("/post", _registerCodeForm).Methods(http.MethodGet)
	r.HandleFunc("/post/register", _registerCode).Methods(http.MethodPost)
	r.HandleFunc("/code/{id:[0-9]+}", _getCode).Methods(http.MethodGet)
	r.HandleFunc("/comment", _registerComment).Methods(http.MethodPost)
	r.HandleFunc("/category", _listCategory).Methods(http.MethodGet)
	r.HandleFunc("/category/{cat:[a-zA-Z0-9_]+}", _getCategoryCode).Methods(http.MethodGet)
	err := http.ListenAndServe(":9927", r)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func die(w http.ResponseWriter) {
	http.Error(w, "bad request", http.StatusBadRequest)
}

func templatePath(name string) string {
	return "/home/joona0825/go/src/alfr.kr/winterstudy/html/" + name
}

func parseFile(name string) (*template.Template, error) {
	return template.ParseFiles(templatePath(name), templatePath("header.html"), templatePath("footer.html"))
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

	t, err := parseFile("index.html")
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
	t, err := parseFile("post.html")
	categories := _getCategoryList()

	if err == nil {
		data := struct {
		 	CaptchaKey string
			Categories []Category
		} {
		 	CaptchaKey: CAPTCHA_KEY,
			Categories: categories,
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
	result := registerPost(r.PostFormValue("captcha"), r.PostFormValue("code"), r.PostFormValue("category"))
	if result == -1 {
		t, _ := parseFile("post_fail.html")
		t.Execute(w, r.PostFormValue("code"))
		return
	}

	t, _ := parseFile("post_success.html")
	t.Execute(w, result)
}

func _getCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err == nil {
		post := getPost(id)
		t, err := parseFile("code.html")
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

func _listCategory(w http.ResponseWriter, r *http.Request) {
	categories := _getCategoryList()
	t, err := parseFile("categoryList.html")

	if err == nil {
		data := struct {
			Categories []Category
		} {
			Categories: categories,
		}
		err := t.Execute(w, data) 
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		log.Println(err.Error())
	}

}

func _getCategoryCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cat := vars["cat"]

	list := getPostListFromCategory(cat)

	t, err := parseFile("category.html")
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