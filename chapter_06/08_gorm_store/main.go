package main

import (
	"fmt"

	"time"

	"gorm.io/driver/sqlite"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Post struct {
	Id        int
	Content   string
	Author    string    `sql:"not null"`
	Comments  []Comment `gorm:"ForeignKey:PostId"`
	CreatedAt time.Time
}

type Comment struct {
	Id        int
	Content   string
	Author    string `sql:"not null"`
	PostId    int
	CreatedAt time.Time
}

var Db *gorm.DB

// connect to the Db
func init() {
	var err error
	Db, err = gorm.Open(sqlite.Open("gwp.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	Db.AutoMigrate(&Post{}, &Comment{})
}

func main() {
	post := Post{Content: "Hello World!", Author: "Sau Sheong"}
	fmt.Println(post)

	Db.Omit("Comments").Create(&post)
	fmt.Println(post)

	comment := Comment{Content: "Good post!", Author: "Joe"}
	Db.Model(&post).Association("Comments").Append(&comment)

	// Get one post
	var readPost Post
	Db.Where("author = ?", "Sau Sheong").First(&readPost)
	fmt.Println(readPost)

	var comments []Comment
	Db.Model(&readPost).Association("Comments").Find(&comments)
	fmt.Println(comments[0])
}
