package main

import "github.com/lotproject/user-service/internal"

func main() {
	app := internal.NewApplication()
	app.Run()

}
