package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"githun.com/farmani/snippetbox/internal/models"
)

// Change the signature of the home handler, so it is defined as a method against
// *application.
func home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			notFound(w)
			return
		}

		snippets, err := App.Snippets.Latest()
		if err != nil {
			serverError(w, err)
			return
		}

		for _, snippet := range snippets {
			fmt.Fprintf(w, "%+v\n", snippet)
		}

		/*w.Header().Add("Link", "</static/css/main.css>; rel=preload; as=style")
		w.Header().Add("Link", "</static/js/main.js>; rel=preload; as=script")
		w.Header().Add("Link", "</static/img/favicon.ico>; rel=preload; as=image/x-icon")

		w.WriteHeader(103)*/

		// Include the navigation partial in the template files.
		/*files := []string{
			"./ui/html/layouts/base.tmpl",
			"./ui/html/layouts/nav.tmpl",
			"./ui/html/pages/home.tmpl",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			// Because the home handler function is now a method against application
			// it can access its fields, including the error logger. We'll write the log
			// message to this instead of the standard logger.
			serverError(w, err)
			return
		}

		err = ts.ExecuteTemplate(w, "base", nil)
		if err != nil {
			// Because the home handler function is now a method against application
			// it can access its fields, including the error logger. We'll write the log
			// message to this instead of the standard logger.
			serverError(w, err)
		}*/
	}
}

// Change the signature of the snippetView handler so it is defined as a method
// against *application.
func snippetView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			notFound(w)
			return
		}

		// Use the SnippetModel object's Get method to retrieve the data for a
		// specific record based on its ID. If no matching record is found,
		// return a 404 Not Found response.
		snippet, err := App.Snippets.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				notFound(w)
			} else {
				serverError(w, err)
			}
			return
		}

		// Initialize a slice containing the paths to the view.tmpl file,
		// plus the base layout and navigation partial that we made earlier.
		files := []string{
			"./ui/html/layouts/base.tmpl",
			"./ui/html/layouts/nav.tmpl",
			"./ui/html/pages/view.tmpl",
		}

		// Parse the template files...
		ts, err := template.ParseFiles(files...)
		if err != nil {
			serverError(w, err)
			return
		}

		// And then execute them. Notice how we are passing in the snippet
		// data (a models.Snippet struct) as the final parameter?
		err = ts.ExecuteTemplate(w, "base", &templateData{Snippet: snippet})
		if err != nil {
			serverError(w, err)
		}
	}
}

// Change the signature of the snippetView handler so it is defined as a method
// against *application.
func snippetCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// Create some variables holding dummy data. We'll remove these later on
		// during the build.
		dto := models.SnippetDto{
			Title:   "O snail",
			Content: "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa",
			Expires: 7,
		}

		// Pass the data to the SnippetModel.Insert() method, receiving the
		// ID of the new record back.
		id, err := App.Snippets.Insert(dto)
		if err != nil {
			serverError(w, err)
			return
		}

		// Redirect the user to the relevant page for the snippet.
		http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
	}
}
