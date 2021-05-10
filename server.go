package main

import (
	"stndalng/db"
	"stndalng/route"

	"github.com/labstack/echo/middleware"
)

func main() {

	db.InitDB()
	e := route.Init()

	e.Use(middleware.BodyLimit("10M"))

	e.Use(middleware.CORS())
	e.Logger.Fatal(e.Start(":4000"))
}
