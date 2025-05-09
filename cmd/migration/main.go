package main

import (
	"context"
	"log"
	"tablelink/src/app"
	"tablelink/src/migration"
)

func main() {
	// Init app context
	ctx := context.TODO()
	if err := app.Init(ctx); err != nil {
		panic(err)
	}

	// Init Migration
	if err := migration.Init(app.GormDB()); err != nil {
		panic(err)
	}

	log.Println("Finish Migration")
}
