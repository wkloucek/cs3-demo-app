package main

import (
	"fmt"

	"github.com/wkloucek/cs3-demo-app/pkg/cs3demoapp"
)

func main() {
	err := cs3demoapp.Start()
	if err != nil {
		fmt.Println(err)
	}
}
