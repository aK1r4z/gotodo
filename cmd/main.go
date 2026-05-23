package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aK1r4z/gotodo/internal/store/postgres"
	"github.com/aK1r4z/gotodo/internal/todo"
)

func main() {
	fmt.Println("Hello, GO!")

	ctx := context.Background()

	pgsql, err := postgres.NewDB(ctx, "")
	if err != nil {
		panic(err)
	}

	todoStore := pgsql.TODOStore()

	todoService, err := todo.NewService(todoStore)
	if err != nil {
		panic(err)
	}

	for i, a := range os.Args {
		fmt.Printf("%d %s\n", i, a)
	}

	switch os.Args[1] {
	case "list":
		str, err := todoService.List(ctx)
		if err != nil {
			panic(err)
		}

		fmt.Println(str)
	case "new":
		title := os.Args[2]
		content := os.Args[3]

		err := todoService.Create(ctx, title, content)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Created %s %s\n", title, content)
	case "sort":
		if err := todoService.Sort(ctx); err != nil {
			panic(err)
		}

		str, err := todoService.List(ctx)
		if err != nil {
			panic(err)
		}

		fmt.Println(str)
	}
}
