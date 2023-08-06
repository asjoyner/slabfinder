package main

import (
	"fmt"

	"github.com/asjoyner/slabfinder/fetcher/stonebasyx"
)

func main() {
	// Load known slabs

	slabs, err := stonebasyx.Fetch()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", slabs)
	// remove known slabs
	// filter slabs by criteria
	// append to known slabs, flush to disk
	// write HTML page of known interesting slabs
	// send notification of new slabs
}
