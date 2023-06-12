package main

import (
	"todo-service/conf"
	"todo-service/src/server"
)

//go:generate go run github.com/99designs/gqlgen generate

func main() {
	if err := conf.Init("config"); err != nil {
		panic(err)
	}

	app := server.NewApp()

	if err := app.Run(); err != nil {
		panic(err)
	}
}
