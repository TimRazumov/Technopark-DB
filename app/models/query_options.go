package models

type QueryString struct {
	Limit int    `query:"limit"`
	Since string `query:"since"`
	Desc  bool   `query:"desc"`
}

func CreateQueryString() QueryString {
	return QueryString{
		Limit: 100,
		Since: "",
		Desc:  false,
	}
}
