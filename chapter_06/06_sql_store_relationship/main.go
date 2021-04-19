package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"

	_ "github.com/mattn/go-sqlite3"
)

type Post struct {
	Id       int
	Content  string
	Author   string
	Comments []Comment
}

type Comment struct {
	Id      int
	Content string
	Author  string
	Post    *Post
}

var Db *sql.DB

func runMigrations() {
	sql, _ := ioutil.ReadFile("setup.sql")
	_, err := Db.Exec(string(sql))
	if err != nil {
		panic(err)
	}
}

// connect to the Db
func init() {
	var err error
	Db, err = sql.Open("sqlite3", "gwp.db")
	if err != nil {
		panic(err)
	}

	runMigrations()
}

// Create a new post
func (post *Post) Create() (err error) {
	sql := "insert into posts (content, author) values ($1, $2) returning id"
	stmt, err := Db.Prepare(sql)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(post.Content, post.Author).Scan(&post.Id)
	return
}

// Create a new comment
func (comment *Comment) Create() (err error) {
	if comment.Post == nil {
		err = errors.New("post not found")
		return
	}

	err = Db.QueryRow("insert into comments (content, author, post_id) values ($1, $2, $3) returning id",
		comment.Content, comment.Author, comment.Post.Id).Scan(&comment.Id)
	return
}

// Get a single post
func GetPost(id int) (post Post, err error) {
	post = Post{}
	err = Db.
		QueryRow("select id, content, author from posts where id = $1", id).
		Scan(&post.Id, &post.Content, &post.Author)

	rows, err := Db.Query("select id, content, author from comments where post_id = $1", id)
	if err != nil {
		return
	}

	for rows.Next() {
		comment := Comment{Post: &post}
		err = rows.Scan(&comment.Id, &comment.Content, &comment.Author)
		if err != nil {
			return
		}
		post.Comments = append(post.Comments, comment)
	}
	rows.Close()

	return
}

func main() {
	post := Post{Content: "Hello World!", Author: "Sau Sheong"}
	post.Create()

	comment := Comment{Content: "Good post!", Author: "Joe", Post: &post}
	comment.Create()

	// Get one post
	readPost, _ := GetPost(post.Id)
	fmt.Println(readPost)
	fmt.Println(readPost.Comments)
	fmt.Println(readPost.Comments[0].Post)
}
