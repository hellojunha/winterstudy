package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Comment struct {
	Id int
	Dt time.Time
	Text string
}

type Post struct {
	Id int
	Dt time.Time
	Code string
	Category []Category // Add relationship in DB
	Comments []Comment
}

// Add same field in DB, related to Post
// Add ["LinkedList", 1], ["boj3986", 1], ["boj10845", 1], ["Fibonacci_Time_Complexity", 1] to DB field before publish service
type Category struct {
	Category string // Unique
	Week int
}

func getDatabase() *sql.DB {
	db, err := sql.Open("mysql", DB_USERNAME + ":" + DB_PASSWORD + "@tcp(127.0.0.1:3306)/" + DB_DATABASE + "?parseTime=true")
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	return db
}

func verifyCaptcha(response string) bool {
	if len(response) == 0 {
		return false
	}

	log.Println("response: " + response)
	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify", url.Values{"secret": {CAPTCHA_SECRET}, "response": {response}})
	if err != nil {
		log.Println(err.Error())
		return false
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	return result["success"].(bool)
}

func getPostList(page int) ([]Post, int) {
	db := getDatabase()
	defer db.Close()

	posts := make([]Post, 0)

	rows, err := db.Query("select id from study_post order by id desc limit ?, ?", page * INDEX_PAGING_NUMBER, INDEX_PAGING_NUMBER)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var id int
			err := rows.Scan(&id)
			if err == nil {
				post := getPost(id)
				if post != nil {
					posts = append(posts, *post)
				} else {
					log.Printf("post %d returned nil", id)
				}
			} else {
				log.Println(err.Error())
			}
		}
	} else {
		log.Println(err.Error())
	}

	log.Printf("returning %d posts", len(posts))

	var count int
	db.QueryRow("select count(*) from study_post").Scan(&count)
	pages := count / INDEX_PAGING_NUMBER
	if count % INDEX_PAGING_NUMBER != 0 {
		pages += 1
	}

	return posts, pages
}

// get Category List (string type slice) from DB 
func _getCategoryList() []Category {
	db := getDatabase()
	defer db.Close()

	cats := make([]Category, 0)

	rows, err := db.Query("select category from study_category") // Check this Query

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&category)
			if err == nil {
				cats = append(cats, &category) 
			} else {
				log.Println(err.Error())
			}
		}
	} else {
		log.Println(err.Error())
	}

	// log.Printf("getting categories")

	return cats
}

// similar as func getPostList() but Query is little different, and not paging  
func getPostListFromCategory(string cat) [] {
	db := getDatabase()
	defer db.Close()

	posts := make([]Post, 0)

	rows, err := db.Query("select id from study_post order by id desc where category = ?", cat) 
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var id int
			err := rows.Scan(&id)
			if err == nil {
				post := getPost(id)
				if post != nil {
					posts = append(posts, *post)
				} else {
					log.Printf("post %d returned nil", id)
				}
			} else {
				log.Println(err.Error())
			}
		}
	} else {
		log.Println(err.Error())
	}

	log.Printf("returning %d posts", len(posts))

	return posts
}

func getPost(id int) *Post {
	db := getDatabase()
	defer db.Close()

	var post Post
	err := db.QueryRow("select id, dt, code from study_post where id = ?", id).Scan(&post.Id, &post.Dt, &post.Code)
	if err != nil {
		log.Println(err)
		return nil
	}

	rows, err := db.Query("select id, dt, text from study_comment where post_id = ?", id)
	if err == nil {
		defer rows.Close()
		comments := make([]Comment, 0)
		for rows.Next() {
			var comment Comment
			err := rows.Scan(&comment.Id, &comment.Dt, &comment.Text)
			if err == nil {
				comments = append(comments, comment)
			} else {
				log.Println(err.Error())
			}
		}

		post.Comments = comments
	} else {
		log.Println(err.Error())
	}

	return &post
}

func registerPost(captchaResp string, code string) int {
	if !verifyCaptcha(captchaResp) {
		return -1
	}

	if len(code) == 0 {
		return -1
	}

	db := getDatabase()
	defer db.Close()

	_, err := db.Exec("insert into study_post (code) values (?)", code)
	if err != nil {
		log.Println(err.Error())
		return -1
	}

	var index int
	err = db.QueryRow("select last_insert_id()").Scan(&index)
	if err == nil {
		return index
	}

	return -1
}

func registerComment(captchaResp string, postId int, text string) bool {
	if !verifyCaptcha(captchaResp) {
		return false
	}

	if len(text) == 0 {
		return false
	}

	db := getDatabase()
	defer db.Close()

	_, err := db.Exec("insert into study_comment (post_id, text) values (?, ?)", postId, text)
	if err != nil {
		log.Println(err.Error())
	}
	return err == nil
}