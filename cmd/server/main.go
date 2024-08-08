package main

import "github.com/datnguyen210/go-muji/internal/routers"

func main() {
	r := routers.NewRouter()
	r.Run(":8080")
}
