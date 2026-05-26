package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/aK1r4z/gotodo/internal/store/postgres"
	"github.com/aK1r4z/gotodo/internal/todo"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Hello, GO!")

	godotenv.Load()

	ctx := context.Background()

	pgsql, err := postgres.NewDB(ctx, os.Getenv("CONNSTR"))
	if err != nil {
		panic(err)
	}

	todoStore := pgsql.TODOStore()

	todoService, err := todo.NewService(todoStore)
	if err != nil {
		panic(err)
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
	default:
		num, err := strconv.Atoi(os.Args[1])
		if err != nil {
			panic(err)
		}

		str, err := todoService.GetByNumber(ctx, uint32(num))
		if err != nil {
			panic(err)
		}

		fmt.Println(str)
	}
}
