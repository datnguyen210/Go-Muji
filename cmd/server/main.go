package main

import (
	"log"
	"net/http"

	"github.com/spf13/viper"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hello from Snippetbox"))
}

func productView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a specific product..."))
}

func productCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new product..."))
}

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

	// viper.SetConfigFile(".env")
	// viper.ReadInConfig()
	// port := viper.GetInt("PORT")

	// r := routers.NewRouter()
	// r.Run(fmt.Sprintf(":%d", port))
}
