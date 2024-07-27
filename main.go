package main

import "github.com/hqdem/go-api-template/cmd"

func main() {
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
