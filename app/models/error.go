package models

type Error struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}
