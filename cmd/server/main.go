package main

import (
	"log"
	"net/http"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	port := ":" + viper.GetString("PORT")

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/product/view", productView)
	mux.HandleFunc("/product/create", productCreate)

	log.Printf("Starting server on %s", port)
	err := http.ListenAndServe(port, mux)
	log.Fatal(err)
}
