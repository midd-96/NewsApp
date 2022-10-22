package main

import (
	"fmt"
	"log"
	"net/http"
	"newsapp/forms"
	"newsapp/models"
	"strconv"

	"github.com/CloudyKit/jet/v6"
	"github.com/go-chi/chi/v5"
)

func (a *application) homeHandler(w http.ResponseWriter, r *http.Request) {

	// u := models.User{
	// 	Email:    "test@gmail.com",
	// 	Password: "password",
	// 	Name:     "Test User",
	// }

	// err := a.Models.Users.Insert(&u)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// a.Models.Posts.Insert("Sample title 1", "http://localhost", u.ID)
	// a.Models.Posts.Insert("Sample title 2", "http://localhost", u.ID)
	// a.Models.Posts.Insert("Sample title 3", "http://localhost", u.ID)
	// a.Models.Posts.Insert("Sample title 4", "http://localhost", u.ID)

	err := r.ParseForm()
	if err != nil {
		a.serverError(w, err)
		return
	}

	filter := models.Filter{
		Query:    r.URL.Query().Get("q"),
		Page:     a.readIntDefault(r, "page", 1),
		PageSize: a.readIntDefault(r, "page_size", 5),
		OrderBy:  r.URL.Query().Get("order_by"),
	}

	posts, meta, err := a.Models.Posts.GetAll(filter)
	if err != nil {
		a.serverError(w, err)
		return
	}

	queryUrl := fmt.Sprintf("page_size=%d&order_by=%s&q=%s", meta.PageSize, filter.OrderBy, filter.Query)
	nextUrl := fmt.Sprintf("%s&page=%d", queryUrl, meta.NextPage)
	prevUrl := fmt.Sprintf("%s&page=%d", queryUrl, meta.PrevPage)

	vars := make(jet.VarMap)
	vars.Set("posts", posts)
	vars.Set("meta", meta)
	vars.Set("nextUrl", nextUrl)
	vars.Set("prevUrl", prevUrl)
	vars.Set("form", forms.New(r.Form))

	err = a.render(w, r, "index", vars)

	if err != nil {
		log.Fatal(err)
	}
}

func (a *application) commentHandler(w http.ResponseWriter, r *http.Request) {

	vars := make(jet.VarMap)

	postId, err := strconv.Atoi(chi.URLParam(r, "postId"))
	if err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	post, err := a.Models.Posts.Get(postId)
	if err != nil {
		a.serverError(w, err)
		return
	}

	comments, err := a.Models.Comments.GetForPost(post.ID)
	if err != nil {
		a.serverError(w, err)
		return
	}

	vars.Set("post", post)
	vars.Set("comments", comments)
	err = a.render(w, r, "comments", vars)
	if err != nil {
		a.serverError(w, err)
		return
	}
}

func (a *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	err := a.render(w, r, "login", nil)
	if err != nil {
		a.serverError(w, err)
		return
	}
}

func (a *application) signupHandler(w http.ResponseWriter, r *http.Request) {

	vars := make(jet.VarMap)
	vars.Set("form", forms.New(r.PostForm))

	err := a.render(w, r, "signup", vars)
	if err != nil {
		a.serverError(w, err)
		return
	}
}

func (a *application) loginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024*2)

	err := r.ParseForm()
	if err != nil {
		a.serverError(w, err)
		return
	}

	form := forms.New(r.PostForm)
	form.Email("email")
	form.MinLength("password", 3)

	if !form.Valid() {
		vars := make(jet.VarMap)
		vars.Set("errors", form.Errors)
		err := a.render(w, r, "login", vars)
		if err != nil {
			a.serverError(w, err)
			return
		}
	}

	user, err := a.Models.Users.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		a.session.Put(r.Context(), "flash", "Login error: "+err.Error())
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	a.session.RenewToken(r.Context())
	a.session.Put(r.Context(), sessionKeyUserId, user.ID)
	a.session.Put(r.Context(), sessionKeyUserName, user.Name)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *application) logoutHandler(w http.ResponseWriter, r *http.Request) {

	a.session.Remove(r.Context(), sessionKeyUserId)
	a.session.Remove(r.Context(), sessionKeyUserName)

	a.session.Destroy(r.Context())
	a.session.RenewToken(r.Context())

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
