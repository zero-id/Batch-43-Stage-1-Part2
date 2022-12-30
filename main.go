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

type Project struct {
	Project_name string
	Start_date   string
	End_date     string
	Description  string
	Technologies []string
}

var Projects = []Project{
	{
		Project_name: "Dumbways Mobile Apps",
		Start_date:   "01-09-2022",
		End_date:     "01-12-2022",
		Description:  "This is Dumbways Mobile Application",
		Technologies: []string{"nextJs", "nodeJs", "reactJs", "typeScript"},
	},
}

func main() {
	route := mux.NewRouter()

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public")))) // /public

	route.HandleFunc("/", helloWorld).Methods("GET")
	route.HandleFunc("/home", home).Methods("GET").Name("home")
	route.HandleFunc("/addMyProject", addMyProject).Methods("GET")
	route.HandleFunc("/home/{id}", projectDetail).Methods("GET")
	route.HandleFunc("/home", addProject).Methods("POST")
	route.HandleFunc("/delet-project/{id}", deletProject).Methods("GET")
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

	respData := map[string]interface{}{
		"Data":     Data,
		"Projects": Projects,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, respData)
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

	project_name := r.PostForm.Get("projectName")
	start_date := r.PostForm.Get("startDate")
	end_date := r.PostForm.Get("endDate")
	description := r.PostForm.Get("des")
	technologies := r.Form["techno"]

	var newProject = Project{
		Project_name: project_name,
		Start_date:   start_date,
		End_date:     end_date,
		Description:  description,
		Technologies: technologies,
	}

	Projects = append(Projects, newProject) // memasukan data newProject ke Projects

	fmt.Println(Projects)

	// fmt.Println("Project_name : " + r.PostForm.Get("projectName"))
	// fmt.Println("Start_date : " + r.PostForm.Get("startDate"))
	// fmt.Println("End_date : " + r.PostForm.Get("endDate"))
	// fmt.Println("Description : " + r.PostForm.Get("des"))
	// fmt.Println("Technologies : ", r.Form["techno"])

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)

}

func deletProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	Projects = append(Projects[:id], Projects[id+1:]...)

	fmt.Println(id)
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
