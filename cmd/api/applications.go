package main

import (
	"fmt"
	"net/http"
)

func (app *application) CreateApplicationHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "add application")
}

func (app *application) GetApplicationsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "get applications")
}

func (app *application) GetApplicationHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "get application")
}

func (app *application) UpdateApplicationHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "update application")
}

func (app *application) DeleteApplicationHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "delete application")
}
