package models

import "github.com/valyala/fasthttp"

type QueryString struct {
	Limit int
	Since string
	Desc  bool
	Sort  string
}

func CreateQueryString(agrs *fasthttp.Args) QueryString {
	var res QueryString
	limit, err := agrs.GetUint("limit")
	if err != nil {
		limit = 100
	}
	res.Limit = limit
	res.Since = string(agrs.Peek("since"))
	res.Desc = agrs.GetBool("desc")
	res.Sort = string(agrs.Peek("sort"))
	if len(res.Sort) == 0 {
		res.Sort = "flat"
	}
	return res
}
