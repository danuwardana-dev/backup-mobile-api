package main

import (
	"backend-mobile-api/app/cmd"
	_ "backend-mobile-api/docs"
)

// @title BEYOND-MOBILE-BACKEND
// @version 1.0
// @description BEYOND-MOBILE-BACKEND
// @in header

// @BasePath /
// @query.collection.format multi
func main() {
	cmd.Execute()
}
