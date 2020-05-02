package models

type Vote struct {
	NickName string `json:"nickname"`
	Voice    int    `json:"voice"`
	Thread   int    `json:"thread,omitempty"`
}

//easyjson:json
type Votes []Vote
