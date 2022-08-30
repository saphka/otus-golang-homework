package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	const original = "Hello, OTUS!"
	fmt.Println(stringutil.Reverse(original))
}
