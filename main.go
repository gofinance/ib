package main

import (
	"./ibtws"
	"fmt"
)

func main() {

	h := ibtws.Make()

	fmt.Printf("Handle = %v\n", h)
}
