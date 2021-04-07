package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
)

type App struct {
	Port string
}

func (a *App) Start() {
	addr := fmt.Sprintf(":%s", a.Port)
	http.Handle("/", logreq(index))
	log.Printf("Starting app on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func env(key, defaultValue string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return val
}

func logreq(f func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("path: %s", r.URL.Path)
		f(w, r)
	})
}

func index(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index.html", struct {
		Name string
	}{
		Name: "Sonic the Hedgehog",
	})
}

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	t, err := template.ParseGlob("templates/*.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error %s", err.Error()), 500)
		return
	}

	err = t.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error %s", err.Error()), 500)
		return
	}
}

func main() {
	server := App{
		Port: env("PORT", "8080"),
	}
	server.Start()
}
