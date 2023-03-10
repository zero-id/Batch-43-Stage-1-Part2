package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"personal-web/connection"
	"personal-web/middleware"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type MetaData struct {
	Title string
	IsLogin bool
	UserName string
	FlashData string
}

var Data = MetaData{
	Title: "Personal Web",
}

type Project struct {
	Id           int
	Project_name string
	Start_date   time.Time
	End_date     time.Time
	Description  string
	Technologies []string
	User_id int
	Duration     string
	Image string
}

type User struct {
    Id       int
    Name     string
    Email    string
    Password string
}

func main() {
	route := mux.NewRouter()

	connection.DatabaseConnection()

	//untuk mengakses Folder
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	route.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads/"))))

	//route html
	route.HandleFunc("/home", home).Methods("GET").Name("home")
	route.HandleFunc("/addMyProject", addMyProject).Methods("GET")
	route.HandleFunc("/project-detail/{id}", projectDetail).Methods("GET")
	route.HandleFunc("/edit-project/{id}", editProject).Methods("GET")
	route.HandleFunc("/addMyProject", middleware.UploadFile(addProject)).Methods("POST")
	route.HandleFunc("/edit-project-input/{id}", middleware.UploadFile(editProjectInput)).Methods("POST")
	route.HandleFunc("/delet-project/{id}", deletProject).Methods("GET")
	route.HandleFunc("/contact", contacMe).Methods("GET")

	//route untuk login dan register
	route.HandleFunc("/login", formLogin).Methods("GET")
	route.HandleFunc("/login", login).Methods("POST")
	route.HandleFunc("/register", formRegister).Methods("GET")
	route.HandleFunc("/register", register).Methods("POST")
	route.HandleFunc("/logout", logout).Methods("GET")

	//untuk menapilkan status running di terminal
	fmt.Println("Server is running on port 8080")

	//tempat untuk menjalankan route
	http.ListenAndServe("localhost:8080", route)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	var result []Project

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
		rows, err := connection.Conn.Query(context.Background(), "SELECT id, project_name, start_date, end_date, description, technologies, image FROM tb_project")

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		for rows.Next() { //next dari pgx dan berfunsi untuk dia akan membaca apa dari connection ketika sudah berhasil menjalankan query
			var each = Project{}

			var err = rows.Scan(&each.Id, &each.Project_name, &each.Start_date, &each.End_date, &each.Description, &each.Technologies, &each.Image)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			each.Duration = selisihDate(each.Start_date, each.End_date)

			result = append(result, each)
		}
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
		user := session.Values["Id"]

		rows, err := connection.Conn.Query(context.Background(), "SELECT tb_project.id, project_name, start_date, end_date, description, technologies, image FROM tb_user LEFT JOIN tb_project ON tb_project.user_id = tb_user.id WHERE tb_project.user_id = $1", user)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		for rows.Next() { 
			var each = Project{}

			var err = rows.Scan(&each.Id, &each.Project_name, &each.Start_date, &each.End_date, &each.Description, &each.Technologies, &each.Image)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			each.Duration = selisihDate(each.Start_date, each.End_date)

			result = append(result, each)
		}
	}

	

	fm := session.Flashes("message")

	var flashes []string
	if len(fm) > 0 {
		session.Save(r, w)
		for _, fl := range fm {
			flashes = append(flashes, fl.(string))
		}
	}

	Data.FlashData = strings.Join(flashes, "")

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

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}

	respData := map[string]interface{}{
		"Data":     Data,
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

	respData := map[string]interface{}{
		"Data":     Data,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, respData)
}

func projectDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/project-detail.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
    session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
        Data.IsLogin = false
    } else {
        Data.IsLogin = session.Values["IsLogin"].(bool)
        Data.UserName = session.Values["Name"].(string)
    }

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	each := Project{}
    err = connection.Conn.QueryRow(context.Background(), "SELECT id, project_name, start_date, end_date, description, technologies, image, user_id FROM tb_project WHERE id=$1", id).Scan(
        &each.Id, &each.Project_name, &each.Start_date, &each.End_date, &each.Description, &each.Technologies, &each.Image, &each.User_id)
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


	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
    session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
        Data.IsLogin = false
    } else {
        Data.IsLogin = session.Values["IsLogin"].(bool)
        Data.UserName = session.Values["Name"].(string)
    }

	user := session.Values["Id"].(int)

	dataContext := r.Context().Value("dataFile")
	image := dataContext.(string)

	// Menghitung Durasi
	startDateTime, _ := time.Parse("2006-01-02", start_date)
	endDateTime, _ := time.Parse("2006-01-02", end_date)


	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_project(project_name, start_date, end_date, description, technologies, image, user_id) VALUES ($1,$2,$3,$4,$5,$6,$7)", project_name, startDateTime, endDateTime, description, technologies, image, user)
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
	err = connection.Conn.QueryRow(context.Background(), "SELECT id, project_name, start_date, end_date, description, technologies, image, user_id FROM tb_project WHERE id=$1 ",id).Scan(&ProjectDetail.Id, &ProjectDetail.Project_name, &ProjectDetail.Start_date, &ProjectDetail.End_date, &ProjectDetail.Description, &ProjectDetail.Technologies, &ProjectDetail.Image, &ProjectDetail.User_id)

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
	
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")	

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}


	// Menghitung Durasi
	startDateTime, _ := time.Parse("2006-01-02", start_date)
	endDateTime, _ := time.Parse("2006-01-02", end_date)

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	dataContext := r.Context().Value("dataFile")
	image := dataContext.(string)

	_, err = connection.Conn.Exec(context.Background(), "UPDATE tb_project SET project_name = $1, start_date = $2, end_date = $3, description = $4, technologies = $5, image = $6 WHERE id = $7 ", project_name, startDateTime, endDateTime, description, technologies, image, id)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("message : " + err.Error()))
        return
    }



	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func formRegister (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/register.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func register (w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
    if err != nil {
        log.Fatal(err)
    }

    name := r.PostForm.Get("name")
    email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_user(name, email, password) VALUES ($1,$2,$3)", name, email, passwordHash)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("message : " + err.Error()))
        return
    }

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")
	session.AddFlash("Successfully register!", "message")
	session.Save(r, w)


	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

func formLogin (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset:utf-8")
	
	var tmpl, err =template.ParseFiles("views/login.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	}

	//cookie/ sebuah sesi yang berupa cookie = storing data/ menyompan data

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID") 

	fm := session.Flashes("message")  //untuk mengambil flash message yang di simpan di session id dengan nama message

	var flashes []string 
	if len(fm) > 0 {
		session.Save(r, w)
		for _, fl := range fm {
			flashes = append(flashes, fl.(string))
		}
	}

	Data.FlashData = strings.Join(flashes, "")

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func login(w http.ResponseWriter, r *http.Request) {
	
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
    session, _ := store.Get(r, "SESSION_ID")


    err := r.ParseForm()
    if err != nil {
        log.Fatal(err)
    }

    email := r.PostForm.Get("email")
    password := r.PostForm.Get("password")

	user := User{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT * FROM tb_user WHERE email=$1", email).Scan(
		&user.Id, 
		&user.Name ,
		&user.Email, 
		&user.Password)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	session.Values["IsLogin"] = true
	session.Values["Name"] = user.Name
	session.Values["Id"] = user.Id
	session.Options.MaxAge = 10800


    session.AddFlash("Login success", "message") 
    session.Save(r, w)

    http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func logout (w http.ResponseWriter, r *http.Request) {
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
    session, _ := store.Get(r, "SESSION_ID")

	session.Options.MaxAge = -1

    session.Save(r, w)

    http.Redirect(w, r, "/home", http.StatusSeeOther)
}
