package main

import (
	"./ibtws"
	"fmt"
)

func main() {
	engine, err := ibtws.Make()
	if err != nil {
		fmt.Printf("error initializing engine: %s\n", err)
		return
	}
	fmt.Printf("engine = %v\n", engine)
}
