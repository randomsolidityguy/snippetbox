package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetbox.randomsolidityguy.net/internal/models"
	"snippetbox.randomsolidityguy.net/internal/validator"

	"github.com/julienschmidt/httprouter" // New import
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
    // Because httprouter matches the "/" path exactly, we can now remove the
    // manual check of r.URL.Path != "/" from this handler.

    snippets, err := app.snippets.Latest()
    if err != nil {
        app.serverError(w, err)
        return
    }

    data := app.newTemplateData(r)
    data.Snippets = snippets

    app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
    // When httprouter is parsing a request, the values of any named parameters
    // will be stored in the request context. We'll talk about request context
    // in detail later in the book, but for now it's enough to know that you can
    // use the ParamsFromContext() function to retrieve a slice containing these
    // parameter names and values like so:
    params := httprouter.ParamsFromContext(r.Context())

    // We can then use the ByName() method to get the value of the "id" named
    // parameter from the slice and validate it as normal.
    id, err := strconv.Atoi(params.ByName("id"))
    if err != nil || id < 1 {
        app.notFound(w)
        return
    }

    snippet, err := app.snippets.Get(id)
    if err != nil {
        if errors.Is(err, models.ErrNoRecord) {
            app.notFound(w)
        } else {
            app.serverError(w, err)
        }
        return
    }

    data := app.newTemplateData(r)
    data.Snippet = snippet

    app.render(w, http.StatusOK, "view.tmpl", data)
}

// Add a new snippetCreate handler, which for now returns a placeholder
// response. We'll update this shortly to show a HTML form.
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
    data := app.newTemplateData(r)
    data.Form = &snippetCreateForm{
        Expires: 365,
    }

    app.render(w, http.StatusOK, "create.tmpl", data)
}

type snippetCreateForm struct {
    Title string `form:"title"`
    Content string `form:"content"`
    Expires int `form:"expires"`
    validator.Validator `form:"-"`
}


// Rename this handler to snippetCreatePost.
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        app.clientError(w, http.StatusBadRequest)
        return
    }

    var form snippetCreateForm

    err = app.decodePostForm(r, &form)
    if err != nil {
        app.clientError(w, http.StatusBadRequest)
        return
    }

    form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
    form.CheckField(validator.MaxCharacters(form.Title, 100), "title", "This field cannot be more than 100 characters long")
    form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
    form.CheckField(validator.PermittedInt(form.Expires, 1,7,365), "expires", "This field must equal 1, 7, 365")
    // Validation

    if !form.Valid() {
        data := app.newTemplateData(r)
        data.Form = form
        app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
        return
    }

    
    id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
    if err != nil {
        app.serverError(w, err)
        return
    }

    http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}