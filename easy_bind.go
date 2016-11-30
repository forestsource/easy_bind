package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type server_info struct {
	ip   string
	port int
}

func main() {
	si := server_info{ip: "127.0.0.1", port: 8080}
	fmt.Println("start server")
	http.HandleFunc("/", home)
	http.HandleFunc("/make", make)
	http.Handle("/statics/", staticFile())
	http.ListenAndServe(si.ip+":"+strconv.Itoa(si.port), nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("statics/home.html")
	err = tmpl.Execute(w, nil)
	if err != nil {
		panic(err)
	}
}

type IPs struct {
	Ip1 string
	Ip2 string
}

func make(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println(r.FormValue("ip1"))
	tmpl, err := template.ParseFiles("statics/config.html")
	tmpv := IPs{"127.0.0.1", "127.0.0.2"}
	err = tmpl.Execute(w, tmpv)
	if err != nil {
		panic(err)
	}
}

func staticFile() http.Handler {
	return http.StripPrefix("/statics/", http.FileServer(http.Dir("statics")))
}
