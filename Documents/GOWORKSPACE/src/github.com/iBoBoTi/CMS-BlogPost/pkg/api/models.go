package api

import (
	"time"
)

type User struct {
	ID        string    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"user_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Posts     []Post
	Followers []*User
}

type Post struct {
	ID          string       `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	PublishedAt time.Time `json:"published_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Author      *User     `json:"author"`
	PostType	string
	Comments    []Comment
	Likes	uint
	UserName string
}

type Comment struct {
	ID string 	`json:"id"`
	Content string `json:"content"`
	Author  *User  `json:"author"`
	PublishedAt time.Time `json:"published_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Post    *Post
	Likes uint
}




