package main

import (
	"fmt"
	"net/http"

	HttpPreloader "github.com/picklejw/go-preloader-http"
)

func main() {
	mux := http.NewServeMux()
	reactAppBuildRoot := "" //"./react-app/build"

	// Create context
	preloader := HttpPreloader.NewHttpPreloaderContext()

	// Register preload routes
	preloader.Get("/", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		fmt.Fprintln(w, `{"name":"joe","email":"joe@example.com", "shouldNotHaveID":"%s"}`, id)
	})
	preloader.Get("/item", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		fmt.Fprintf(w, `{"queryParam":"%s","name":"Snowcone"}`, id)
	})

	// Wrap with middleware
	handler := preloader.HttpPreloader(mux, "/api", reactAppBuildRoot)

	fmt.Println("Server listening on :8888")
	http.ListenAndServe(":8888", handler)
}
