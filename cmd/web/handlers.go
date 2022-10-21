package main

import (
	"fmt"
	"log"
	"net/http"
	"newsapp/models"

	"github.com/CloudyKit/jet/v6"
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
	//vars.Set("form", forms.New(r.Form))

	err = a.render(w, r, "index", vars)

	if err != nil {
		log.Fatal(err)
	}
}
