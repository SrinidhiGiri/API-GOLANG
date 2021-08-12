package main

import (
	"awesomeProject6/model"
	db "awesomeProject6/repository"
	. "awesomeProject6/router"
	"fmt"
)

func main() {
	dbHost := "127.0.0.1:27017"
	fmt.Println("Connected to MongoDB")
	db.Init(&model.Database{
		Driver:   "mongodb",
		Endpoint: dbHost})
	defer db.Exit()

	Router()

}