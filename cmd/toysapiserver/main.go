package main

import (
	"fmt"

	"github.com/lucku/otto-coding-challenge/api/mytoystestapi"
)

/* Start server, nothing more */

func main() {

	test := mytoystestapi.GetAllLinks()

	fmt.Println(test)
}
