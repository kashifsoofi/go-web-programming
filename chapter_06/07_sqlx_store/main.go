package main

import (
	"fmt"
	"io/ioutil"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Post struct {
	Id      int
	Content string
	Author  string
}

var Db *sqlx.DB

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
	Db, err = sqlx.Open("sqlite3", "gwp.db")
	if err != nil {
		panic(err)
	}

	runMigrations()
}

// Create a new post
func (post *Post) Create() (err error) {
	sql := "insert into posts (content, author) values ($1, $2) returning id"
	err = Db.QueryRow(sql, post.Content, post.Author).Scan(&post.Id)
	return
}

// Get a single post
func GetPost(id int) (post Post, err error) {
	post = Post{}
	err = Db.
		QueryRowx("select id, content, author from posts where id = $1", id).
		StructScan(&post)
	return
}

func main() {
	post := Post{Content: "Hello World!", Author: "Sau Sheong"}
	post.Create()

	// Get one post
	readPost, _ := GetPost(post.Id)
	fmt.Println(readPost)
}
