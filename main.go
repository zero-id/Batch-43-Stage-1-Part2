package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

var Data = map[string]interface{}{
	"Title": "Personal Web",
}

type Project struct {
	Id           int
	Project_name string
	Start_date   string
	End_date     string
	Description  string
	Technologies []string
	Duration     string
}

var Projects = []Project{
	// {
	// 	Project_name: "Dumbways Mobile Apps",
	// 	Start_date:   "01-09-2022",
	// 	End_date:     "01-12-2022",
	// 	Description:  "This is Dumbways Mobile Application",
	// 	Technologies: []string{"nextJs", "nodeJs", "reactJs", "typeScript"},
	// },
}

func main() {
	route := mux.NewRouter()

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public")))) // /public

	route.HandleFunc("/", helloWorld).Methods("GET")
	route.HandleFunc("/home", home).Methods("GET").Name("home")
	route.HandleFunc("/addMyProject", addMyProject).Methods("GET")
	route.HandleFunc("/project-detail/{id}", projectDetail).Methods("GET")
	route.HandleFunc("/home", addProject).Methods("POST")
	route.HandleFunc("/edit-project/{id}", editProject).Methods("GET")
	route.HandleFunc("/edit-project-input/{id}", editProjectInput).Methods("POST")
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

	var tmpl, err = template.ParseFiles("views/project-detail.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	ProjectDetail := Project{}

	for index, data := range Projects {
		if index == id {
			newStartDate, _ := time.Parse("2006-01-02", data.Start_date)
			newEndDate, _ := time.Parse("2006-01-02", data.End_date)

			ProjectDetail = Project{
				Id:           id,
				Project_name: data.Project_name,
				Start_date:   newStartDate.Format("02 Jan 2006"),
				End_date:     newEndDate.Format("02 Jan 2006"),
				Description:  data.Description,
				Technologies: data.Technologies,
				Duration:     data.Duration,
			}
		}
	}

	resp := map[string]interface{}{
		"Data":          Data,
		"ProjectDetail": ProjectDetail,
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

	// Menghitung Durasi
	startDateTime, _ := time.Parse("2006-01-02", start_date)
	endDateTime, _ := time.Parse("2006-01-02", end_date)

	selisihDate := endDateTime.Sub(startDateTime)

	var duration string
	year := int(selisihDate.Hours() / (12 * 30 * 24))
	if year != 0 {
		duration = strconv.Itoa(year) + " tahun"
	} else {
		month := int(selisihDate.Hours() / (30 * 24))
		if month != 0 {
			duration = strconv.Itoa(month) + " bulan"
		} else {
			week := int(selisihDate.Hours() / (7 * 24))
			if week != 0 {
				duration = strconv.Itoa(week) + " minggu"
			} else {
				day := int(selisihDate.Hours() / (24))
				if day != 0 {
					duration = strconv.Itoa(day) + " hari"
				}
			}
		}
	}

	var newProject = Project{
		Project_name: project_name,
		Start_date:   start_date,
		End_date:     end_date,
		Description:  description,
		Technologies: technologies,
		Duration:     duration,
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

func editProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/edit-project.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	ProjectDetail := Project{}

	for index, data := range Projects {
		if index == id {
			ProjectDetail = Project{
				Id:           id,
				Project_name: data.Project_name,
				Start_date:   data.Start_date,
				End_date:     data.End_date,
				Description:  data.Description,
				Technologies: data.Technologies,
				Duration:     data.Duration,
			}
		}
	}

	resp := map[string]interface{}{
		"Data":          Data,
		"ProjectDetail": ProjectDetail,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)
}

func editProjectInput(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	project_name := r.PostForm.Get("projectName")
	start_date := r.PostForm.Get("startDate")
	end_date := r.PostForm.Get("endDate")
	description := r.PostForm.Get("des")
	technologies := r.Form["techno"]

	// Menghitung Durasi
	startDateTime, _ := time.Parse("2006-01-02", start_date)
	endDateTime, _ := time.Parse("2006-01-02", end_date)

	selisihDate := endDateTime.Sub(startDateTime)

	var duration string
	year := int(selisihDate.Hours() / (12 * 30 * 24))
	if year != 0 {
		duration = strconv.Itoa(year) + " tahun"
	} else {
		month := int(selisihDate.Hours() / (30 * 24))
		if month != 0 {
			duration = strconv.Itoa(month) + " bulan"
		} else {
			week := int(selisihDate.Hours() / (7 * 24))
			if week != 0 {
				duration = strconv.Itoa(week) + " minggu"
			} else {
				day := int(selisihDate.Hours() / (24))
				if day != 0 {
					duration = strconv.Itoa(day) + " hari"
				}
			}
		}
	}

	var newProject = Project{
		Project_name: project_name,
		Start_date:   start_date,
		End_date:     end_date,
		Description:  description,
		Technologies: technologies,
		Duration:     duration,
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	Projects[id] = newProject

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
