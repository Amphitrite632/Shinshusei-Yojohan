package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"encoding/json"
	"strconv"
	"text/template"
	"time"
)

func topPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		notFoundHandler(w, r, http.StatusNotFound)
		return
	}
	template := template.Must(template.ParseFiles("index.html"))
	article := map[string]string{
		"now": fmt.Sprint(time.Now().Unix()),
	}
	template.ExecuteTemplate(w, "index.html", article)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	http.ServeFile(w, r, "resource/404.html")
}

func articleRequestHandler(w http.ResponseWriter, r *http.Request) {
	template := template.Must(template.ParseFiles("articles/template.html"))

	var paths []string
	jsonData, _ := ioutil.ReadFile("./articles/index.json")
	json.Unmarshal(jsonData, &paths)

	maxPageNumber := math.Ceil(float64(len(paths)) / 3.0)
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))

	if page > int(maxPageNumber) || page < 1 {
		page = 1
	}

	var text string
	for i := page*3 - 3; i < (page-1)*3+3; i++ {
		if i < len(paths) {
			t, _ := ioutil.ReadFile("articles/" + paths[i])
			text += "<div class= \"line\"></div>" + string(t)
		}
	}

	var navi string
	if page == 1 {
		navi = "<p class=\"navi\"><a href=\"/articles?page=1\">1</a><a href=\"/articles?page=2\">2</a><a href=\"/articles?page=3\">3</a><a>...</a><a href=\"/articles?page=" + fmt.Sprint(int(maxPageNumber)) + "\">" + fmt.Sprint(int(maxPageNumber)) + "</a></p> "
	} else if page == int(maxPageNumber) {
		navi = "<p class=\"navi\"><a href=\"/articles?page=1\">1</a><a>...</a><a href=\"/articles?page=" + fmt.Sprint(page-2) + "\">" + fmt.Sprint(page-2) + "</a><a href=\"/articles?page=" + fmt.Sprint(page-1) + "\">" + fmt.Sprint(page-1) + "</a><a href=\"/articles?page=" + fmt.Sprint(page) + "\">" + fmt.Sprint(page) + "</a>"
	} else {
		navi = "<p class=\"navi\"><a href=\"/articles?page=1\">1</a><a>...</a><a href=\"/articles?page=" + fmt.Sprint(page-1) + "\">" + fmt.Sprint(page-1) + "</a><a href=\"/articles?page=" + fmt.Sprint(page) + "\">" + fmt.Sprint(page) + "</a><a href=\"/articles?page=" + fmt.Sprint(page+1) + "\">" + fmt.Sprint(page+1) + "</a><a>...</a><a href=\"/articles?page=" + fmt.Sprint(int(maxPageNumber)) + "\">" + fmt.Sprint(int(maxPageNumber)) + "</a></p> "
	}

	article := map[string]string{
		"notes": string(text),
		"navi":  navi,
	}

	template.ExecuteTemplate(w, "template.html", article)
}

func linksRequestHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "links/index.html")
}

func main() {
	http.Handle("/resource/", http.StripPrefix("/resource/", http.FileServer(http.Dir("resource/"))))
	http.HandleFunc("/articles/", articleRequestHandler)
	http.HandleFunc("/links/", linksRequestHandler)
	http.HandleFunc("/", topPageHandler)
	http.ListenAndServe(":8500", nil)
}
