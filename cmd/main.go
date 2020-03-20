package main

import (
	"fmt"

	"github.com/TimRazumov/Technopark-DB/app/test"

    "github.com/rs/xid"
)

func main() {
	test.Hello()

	fmt.Println(xid.New().String())
}
