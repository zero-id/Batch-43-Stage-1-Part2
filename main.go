package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"personal-web/connection"
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
	Start_date   time.Time
	End_date     time.Time
	Description  string
	Technologies []string
	Duration     string
	// Image string
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

	connection.DatabaseConnection()

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public")))) // /public

	route.HandleFunc("/", helloWorld).Methods("GET")
	route.HandleFunc("/home", home).Methods("GET").Name("home")
	route.HandleFunc("/addMyProject", addMyProject).Methods("GET")
	route.HandleFunc("/project-detail/{id}", projectDetail).Methods("GET")
	route.HandleFunc("/addMyProject", addProject).Methods("POST")
	route.HandleFunc("/edit-project/{id}", editProject).Methods("GET")
	route.HandleFunc("/edit-project-input/{id}", editProjectInput).Methods("POST")
	route.HandleFunc("/delet-project/{id}", deletProject).Methods("GET")
	route.HandleFunc("/contact", contacMe).Methods("GET")
	fmt.Println("Server is running on port 5000")
	http.ListenAndServe("localhost:8080", route)
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

	rows, _ := connection.Conn.Query(context.Background(), "SELECT * FROM tb_project")

	var result []Project
	for rows.Next() { //next dari pgx dan berfunsi untuk dia akan membaca apa dari connection ketika sudah berhasil menjalankan query, berarti next akan membaca valuenya yang di kirimkan database
		var each = Project{}

		var err = rows.Scan(&each.Id, &each.Project_name, &each.Start_date, &each.End_date, &each.Description, &each.Technologies)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		each.Duration = selisihDate(each.Start_date, each.End_date)

		result = append(result, each)
	}

	fmt.Println(result)

	respData := map[string]interface{}{
		"Data":     Data,
		"Projects": result,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, respData)
}

func selisihDate(start time.Time, end time.Time ) string {
	selisihDate := end.Sub(start)

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


		return duration 
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

	each := Project{}
    err = connection.Conn.QueryRow(context.Background(), "SELECT id, project_name, start_date, end_date, description, technologies FROM tb_project WHERE id=$1", id).Scan(
        &each.Id, &each.Project_name, &each.Start_date, &each.End_date, &each.Description, &each.Technologies)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("message : " + err.Error()))
        return
    }

	each.Duration = selisihDate(each.Start_date, each.End_date)


	resp := map[string]interface{}{
		"Data":          Data,
		"ProjectDetail": each,
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


	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_project(project_name, start_date, end_date, description, technologies) VALUES ($1,$2,$3,$4,$5)", project_name, startDateTime, endDateTime, description, technologies)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("message : " + err.Error()))
        return
    }

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)

}

func deletProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_project WHERE id=$1", id)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("message : " + err.Error()))
        return
    }

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
	err = connection.Conn.QueryRow(context.Background(), "SELECT id, project_name, start_date, end_date, description, technologies FROM tb_project WHERE id=$1 ",id).Scan(&ProjectDetail.Id, &ProjectDetail.Project_name, &ProjectDetail.Start_date, &ProjectDetail.End_date, &ProjectDetail.Description, &ProjectDetail.Technologies)

    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("message : " + err.Error()))
        return
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

	

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	_, err = connection.Conn.Exec(context.Background(), "UPDATE tb_project SET project_name = $1, start_date = $2, end_date = $3, description = $4, technologies = $5 WHERE id = $6 ", project_name, startDateTime, endDateTime, description, technologies, id)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("message : " + err.Error()))
        return
    }

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}


