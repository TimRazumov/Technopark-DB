package models

import (
	"github.com/valyala/fasthttp"
	"strings"
	"time"
)

type Post struct {
	ID       int       `json:"id,omitempty"`
	Parent   int       `json:"parent,omitempty"`
	Author   string    `json:"author"`
	Message  string    `json:"message"`
	IsEdited bool      `json:"isEdited,omitempty"`
	Forum    string    `json:"forum,omitempty"`
	Thread   int       `json:"thread,omitempty"`
	Created  time.Time `json:"created,omitempty"`
	Path     []int64   `json:"-"`
}

//easyjson:json
type Posts []Post

type Related struct {
	User   bool
	Forum  bool
	Thread bool
}

func CreateRelated(agrs *fasthttp.Args) Related {
	var res Related
	query := strings.Split(string(agrs.Peek("related")), ",")
	for _, param := range query {
		if param == "user" {
			res.User = true
			continue
		}
		if param == "forum" {
			res.Forum = true
			continue
		}
		if param == "thread" {
			res.Thread = true
			continue
		}
	}
	return res
}

type PostGet struct {
	Post   *Post   `json:"post"`
	Author *User   `json:"author,omitempty"`
	Thread *Thread `json:"thread,omitempty"`
	Forum  *Forum  `json:"forum,omitempty"`
}

//easyjson:json
type PostGets []PostGet
