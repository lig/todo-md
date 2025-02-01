package main

import (
	"log"

	"codeberg.org/lig/todo-md/internal/app/todomd"
)

func main() {
	err := todomd.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
