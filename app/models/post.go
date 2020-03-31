package models

import "time"

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

type Related struct {
	User   bool
	Forum  bool
	Thread bool
}

type PostGet struct {
	Post   *Post   `json:"post"`
	Author *User   `json:"author,omitempty"`
	Thread *Thread `json:"thread,omitempty"`
	Forum  *Forum  `json:"forum,omitempty"`
}
