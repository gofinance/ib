package main

import (
	"./ibtws"
	"fmt"
)

func main() {
	engine, err := ibtws.Make()
    if err != nil {
        fmt.Printf("error %v initializing engine\n", err)
        return
    }
	fmt.Printf("engine = %v\n", engine)
    for {
        packet := <-engine.Receive()
        fmt.Printf("packet = %v\n", packet)
    }
}
