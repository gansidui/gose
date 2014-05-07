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
		q := r.Form.Get("q")
		if q == "" {
			utils.RespondTemplate(w, http.StatusOK, "views/html/index.html", nil)
		} else {
			s := "/search?q=" + q + "&start=0"
			http.Redirect(w, r, s, 303)
		}
	} else {
		utils.RespondTemplate(w, http.StatusOK, "views/html/index.html", nil)
	}
}

func SearchPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		q := r.Form.Get("q")
		s := "/search?q=" + q + "&start=0"
		http.Redirect(w, r, s, 303)
	} else {
		q := r.URL.Query().Get("q")
		start := r.URL.Query().Get("start")
		num := "10"
		resultPage, success := dao.GetResultPageInfo(q, start, num)
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
