// Copyright © 2022 UCloud. All rights reserved.

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Application struct {
}

func (app *Application) describe(w http.ResponseWriter, r *http.Request) {
	fmt.Printf(
		"request: %v\n",
		r.Header.Get("Content-Type"),
	)
	fmt.Fprintln(w, "/describe", r.URL.EscapedPath())
}

func (app *Application) create(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("request: %v\n", r)
	fmt.Fprintln(w, "create", r.URL.EscapedPath())
}

func (app *Application) start(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("request: %v\n", r)
	fmt.Fprintln(w, "start", r.URL.EscapedPath())
}

func (app *Application) Routes() http.Handler {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/describe", app.describe)
	router.HandlerFunc(http.MethodPost, "/create", app.create)
	router.HandlerFunc(http.MethodGet, "/start", app.start)

	return router
}

func main() {
	app := Application{}
	srv := &http.Server{
		Addr:         ":8089",
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln("errors: ", err)
	}
}
