package main

import (
	"fmt"
	"golang.org/x/example/stringutil"
)

func main() {
	const original = "Hello, OTUS!"
	var reversed = stringutil.Reverse(original)
	fmt.Println(reversed)
}
