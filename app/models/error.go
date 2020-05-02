package models

import (
	"net/http"
	"strconv"
)

type Error struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

func (err *Error) GetMessage() []byte {
	mess, _ := err.MarshalJSON()
	return mess
}

func CreateNotFoundUser(nickName string) *Error {
	return &Error{
		Code:    http.StatusNotFound,
		Message: "Can't find user by nickname: " + nickName,
	}
}

func CreateConflictUser(nickName string) *Error {
	return &Error{
		Code:    http.StatusConflict,
		Message: "This email is already registered by user: " + nickName,
	}
}

func CreateNotFoundForum(slug string) *Error {
	return &Error{
		Code:    http.StatusNotFound,
		Message: "Can't find forum with slug: " + slug,
	}
}

func CreateNotFoundAuthorThread(nickName string) *Error {
	return &Error{
		Code:    http.StatusNotFound,
		Message: "Can't find thread author by nickname: " + nickName,
	}
}

func CreateNotFoundForumThread(slug string) *Error {
	return &Error{
		Code:    http.StatusNotFound,
		Message: "Can't find thread forum by slug: " + slug,
	}
}

func CreateConflictPost() *Error {
	return &Error{
		Code:    http.StatusConflict,
		Message: "Parent post was created in another thread",
	}
}

func CreateNotFoundAuthorPost(nickName string) *Error {
	return &Error{
		Code:    http.StatusNotFound,
		Message: "Can't find thread author by nickname: " + nickName,
	}
}

func CreateNotFoundThreadPost(id int) *Error {
	return &Error{
		Code:    http.StatusNotFound,
		Message: "Can't find post thread by id: " + strconv.Itoa(id),
	}
}
