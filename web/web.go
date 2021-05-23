package web

import (
	"../api"
	"log"
	"net/http"
	"html/template"

)

func Web() {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/max", MaxService)
	http.HandleFunc("/min", MinService)
	http.HandleFunc("/rand", RandService)
	http.HandleFunc("/choise", ChoiseService)
	http.HandleFunc("/adminCheck", AdminCheck)

	err := http.ListenAndServe("localhost:8080", nil)
	log.Fatal(err)
}


func ChoiseService(writer http.ResponseWriter, request *http.Request) {
	name := request.URL.Query().Get("name")
	html, err := template.ParseFiles("../../web/template/choise.html")
	checkErr(err)
	err = html.Execute(writer, api.ChoiseRow(name))
	checkErr(err)
}

func MaxService(writer http.ResponseWriter, request *http.Request) {
	html, err := template.ParseFiles("../../web/template/max.html")
	checkErr(err)
	err = html.Execute(writer, api.MaxTimeRow())
	checkErr(err)
}
 
func MinService(writer http.ResponseWriter, request *http.Request) {
	html, err := template.ParseFiles("../../web/template/min.html")
	checkErr(err)
	err = html.Execute(writer, api.MinTimeRow())
	checkErr(err)
}

func RandService(writer http.ResponseWriter, request *http.Request) {
	html, err := template.ParseFiles("../../web/template/rand.html")
	checkErr(err)
	err = html.Execute(writer, api.RandomRow())
	checkErr(err)
}

func IndexHandler(writer http.ResponseWriter, request *http.Request) {
	html, err := template.ParseFiles("../../web/template/index.html")
	checkErr(err)
	err = html.Execute(writer, map[string][]api.Services{"Services": api.GetDomains()})
	checkErr(err)
}

func AdminCheck(writer http.ResponseWriter, request *http.Request) {
	html, err := template.ParseFiles("../../web/template/adminCheck.html")
	checkErr(err)
	err = html.Execute(writer, map[string][]api.Services{"Services": api.GetDomains()})
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
