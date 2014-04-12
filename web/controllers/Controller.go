package controllers

import (
	"github.com/gansidui/gose/web/models/dao"
	"github.com/gansidui/gose/web/utils"
	"html/template"
	"log"
	"net/http"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		word := r.Form.Get("word")
		if word == "" {
			utils.RespondTemplate(w, http.StatusOK, "views/html/index.html", nil)
		} else {
			s := "/search?word=" + word + "&page=1"
			http.Redirect(w, r, s, 303)
		}
	} else {
		utils.RespondTemplate(w, http.StatusOK, "views/html/index.html", nil)
	}
}

func SearchPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		word := r.Form.Get("word")
		s := "/search?word=" + word + "&page=1"
		http.Redirect(w, r, s, 303)
	} else {
		word := r.URL.Query().Get("word")
		page := r.URL.Query().Get("page")
		resultPage, success := dao.GetResultPageInfo(word, page)
		if success {
			w.WriteHeader(http.StatusOK)
			t, err := template.ParseFiles("views/html/search.html", "views/html/pagination.html")
			if err != nil {
				log.Fatal(err)
			}
			t.Execute(w, &resultPage)
		} else {
			utils.RespondNotFound(w)
		}
	}
}
