package main

import (
	"fmt"

	"github.com/frigus02/website/generator/data"
)

func main() {
	posts, err := data.ReadProjects()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Posts: %v\n", posts)
}
