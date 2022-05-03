package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

const phrase = "Hello, OTUS!"

func main() {
	fmt.Println(stringutil.Reverse(phrase))
}
