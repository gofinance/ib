package algokit

import (
	"fmt"
	"log"
)

// func Log(v ...): loging. give log information if debug is true

func trace(v ...interface{}) {
	ret := fmt.Sprint(v)
	log.Printf("CLIENT: %s", ret)
}

// func test(): testing for error

func test(err error, mesg string) {
	if err != nil {
		log.Fatalf("CLIENT: ERROR: %v when %s", err, mesg)
	} else {
		trace("Ok: ", mesg)
	}
}
