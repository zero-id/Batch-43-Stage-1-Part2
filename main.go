package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
)

var Data = map[string]interface{}{
	"Title": "Personal Web",
}

func main() {
	route := mux.NewRouter()

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public")))) // /public

	route.HandleFunc("/", helloWorld).Methods("GET")
	route.HandleFunc("/home", home).Methods("GET").Name("home")
	route.HandleFunc("/addMyProject", addMyProject).Methods("GET")
	route.HandleFunc("/home/{id}", projectDetail).Methods("GET")
	route.HandleFunc("/home", addProject).Methods("POST")
	route.HandleFunc("/contact", contacMe).Methods("GET")

	fmt.Println("Server is running on port 5000")
	http.ListenAndServe("localhost:5000", route)
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello world!"))
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func addMyProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/my-project.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func projectDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var tmpl, err = template.ParseFiles("views/project-detail.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	resp := map[string]interface{}{
		"Data": Data,
		"Id":   id,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)
}

func addProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Project Name : " + r.PostForm.Get("projectName"))
	fmt.Println("Start Date : " + r.PostForm.Get("startDate"))
	fmt.Println("End Date : " + r.PostForm.Get("endDate"))
	fmt.Println("Description : " + r.PostForm.Get("des"))
	fmt.Println("Technologies : " + r.PostForm.Get("techno"))

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)

}

func contacMe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Contact-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/contatc-form.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}
