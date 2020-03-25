package models

type User struct {
	NickName string `json:"nickname,omitempty"`
	FullName string `json:"fullname"`
	Email    string `json:"email"`
	About    string `json:"about,omitempty"`
}
