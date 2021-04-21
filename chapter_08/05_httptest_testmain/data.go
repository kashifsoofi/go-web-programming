package main

import (
	"database/sql"
	"io/ioutil"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

// connect to the Db
func init() {
	var err error
	Db, err = sql.Open("sqlite3", "gwp.db")
	if err != nil {
		panic(err)
	}

	runMigrations()
}

func runMigrations() {
	sql, _ := ioutil.ReadFile("setup.sql")
	_, err := Db.Exec(string(sql))
	if err != nil {
		panic(err)
	}
}

// Get a single post
func retrieve(id int) (post Post, err error) {
	post = Post{}
	err = Db.QueryRow("select id, content, author from posts where id = $1", id).Scan(&post.Id, &post.Content, &post.Author)
	return
}

// Create a new post
func (post *Post) create() (err error) {
	statement := "insert into posts (content, author) values ($1, $2) returning id"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(post.Content, post.Author).Scan(&post.Id)
	return
}

// Update a post
func (post *Post) update() (err error) {
	_, err = Db.Exec("update posts set content = $2, author = $3 where id = $1", post.Id, post.Content, post.Author)
	return
}

// Delete a post
func (post *Post) delete() (err error) {
	_, err = Db.Exec("delete from posts where id = $1", post.Id)
	return
}
