package api

import (
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"net/http"
)

func Ready() {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/test", GetTest),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)

	http.Handle("/api/v1/", http.StripPrefix("/api/v1", api.MakeHandler()))
	log.Fatal(http.ListenAndServe(":6000", nil))
}

func GetTest(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(map[string]string{"Body": "Hello World!"})
}
