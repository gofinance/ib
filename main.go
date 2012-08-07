package main

import (
	"./trade"
	"fmt"
)

func main() {
	engine, err := trade.Make()
	if err != nil {
		fmt.Printf("error initializing engine: %s\n", err)
		return
	}
	fmt.Printf("engine = %v\n", engine)
}
