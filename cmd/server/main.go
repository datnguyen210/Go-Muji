package main

import (
	"fmt"

	"github.com/datnguyen210/go-muji/internal/routers"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	port := viper.GetInt("PORT")

	r := routers.NewRouter()
	r.Run(fmt.Sprintf(":%d", port))
}
