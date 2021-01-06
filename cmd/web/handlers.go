package main

import (
	"errors"
	"fmt"
	"net/http"
	"softarch/pkg/models"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	for _, snippet := range s {
		fmt.Fprintf(w, "%v\n", snippet)
	}

	//files := []string {
	//	"./ui/html/home.page.tmpl",
	//	"./ui/html/base.layout.tmpl",
	//	"./ui/html/footer.partial.tmpl",
	//}
	//
	//ts,err := template.ParseFiles(files...)
	//if err!=nil {
	//	app.serverError(w, err)
	//	http.Error(w, "Internal server error", http.StatusInternalServerError)
	//}
	//err = ts.Execute(w, nil)
	//if err!=nil {
	//	app.serverError(w, err)
	//	http.Error(w, "Internal server error", http.StatusInternalServerError)
	//}
}

func (app *application) ShowSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 0 {
		app.notFound(w)
		return
	}
	//TODO
	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	fmt.Fprintf(w, "%v", s)

	//files := []string {
	//	"./ui/html/home.page.tmpl",
	//	"./ui/html/base.layout.tmpl",
	//	"./ui/html/footer.partial.tmpl",
	//}
	//
	//ts,err := template.ParseFiles(files...)
	//if err!=nil {
	//	app.serverError(w, err)
	//	http.Error(w, "Internal server error", http.StatusInternalServerError)
	//}
	//err = ts.Execute(w, nil)
	//if err!=nil {
	//	app.serverError(w, err)
	//	http.Error(w, "Internal server error", http.StatusInternalServerError)
	//}

}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "0 snail"
	content := "0 snail\nClimb Mount Fuji,\\nBut slowly, slowly!\\n\\n– Kobayashi Issa"
	expires := "7"

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
