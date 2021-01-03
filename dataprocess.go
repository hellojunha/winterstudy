package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"net/url"
)

type Comment struct {
	id int
	ts float64
	text string
}

type Post struct {
	id int
	ts float64
	code string
	comments []Comment
}

func getDatabase() *sql.DB {
	db, err := sql.Open("mysql", DB_USERNAME + ":" + DB_PASSWORD + "@tcp(127.0.0.1:3306)/" + DB_DATABASE)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	return db
}

func verifyCaptcha(response string) bool {
	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify", url.Values{"secret": {""}, "response": {response}})
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

func getPostList(page int) []Post {
	db := getDatabase()
	defer db.Close()

	posts := make([]Post, 0)

	rows, err := db.Query("select id from study_post order by ts desc limit ?, 10", page * 10)
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
				log.Println(err)
			}
		}
	}

	return posts
}

func getPost(id int) *Post {
	db := getDatabase()
	defer db.Close()

	var post Post
	err := db.QueryRow("select id, ts, code from study_post where id = ?", id).Scan(post.id, post.ts, post.code)
	if err != nil {
		log.Println(err)
		return nil
	}

	rows, err := db.Query("select id, ts, text from study_comments where post_id = ?", id)
	if err == nil {
		defer rows.Close()
		comments := make([]Comment, 0)
		for rows.Next() {
			var comment Comment
			err := rows.Scan(&comment.id, &comment.ts, &comment.text)
			if err == nil {
				comments = append(comments, comment)
			} else {
				log.Println(err)
			}
		}

		post.comments = comments
	} else {
		log.Println(err)
	}

	return &post
}

func registerPost(captchaResp string, code string) bool {
	if !verifyCaptcha(captchaResp) {
		return false
	}

	db := getDatabase()
	defer db.Close()

	_, err := db.Exec("insert into study_post (code) values (?)", code)
	if err != nil {
		log.Println(err.Error())
	}
	return err == nil
}

func registerComment(captchaResp string, postId int, text string) bool {
	if !verifyCaptcha(captchaResp) {
		return false
	}

	db := getDatabase()
	defer db.Close()

	_, err := db.Exec("insert into study_comments (post_id, text) values (?, ?)", postId, text)
	if err != nil {
		log.Println(err.Error())
	}
	return err == nil
}