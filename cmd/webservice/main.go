package main

import (
	"os"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var app App
var config AuthConfig

func init() {
	ct, err := app.GetDBConnectionType()
	if err != nil {
		panic(err)
	}

	config = AuthConfig{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	}

	err = app.Initialize(ct, config)
	if err != nil {
		panic(err)
	}
}

func main() {
	defer func() {
		err := app.DB.Close()
		if err != nil {
			panic(err)
		}
	}()
	app.Run("8080")
}
