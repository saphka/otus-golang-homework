package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	const original = "Hello, OTUS!"
	reversed := stringutil.Reverse(original)
	fmt.Println(reversed)
}
