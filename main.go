package main

import (
	"fmt"
	"mt-hosting-manager/db"
	"mt-hosting-manager/web"
	"os"
)

func main() {
	fmt.Println("Listening on port 8080")

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	db_, err := db.Init(wd)
	if err != nil {
		panic(err)
	}

	err = db.Migrate(db_)
	if err != nil {
		panic(err)
	}

	repos := db.NewRepositories(db_)
	err = web.Serve(repos)
	if err != nil {
		panic(err)
	}
}
