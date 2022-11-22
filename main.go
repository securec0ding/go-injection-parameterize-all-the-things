package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type EmployeePageData struct {
	PageTitle string
	Employees []Employee
}

type Employee struct {
	Name   string
	Email  string
	Phone  string
	DOB    string
	Salary int
}

func getEmployees(search string) []Employee {
	db, err := sql.Open("sqlite3", "/tmp/employees.db")
	if err != nil {
		panic(err)
	}

	rows, _ := db.Query("SELECT * FROM employee WHERE name LIKE '%" + search + "%';")

	var results []Employee

	var id int
	var name string
	var email string
	var phone string
	var dob string
	var salary int
	for rows.Next() {
		_ = rows.Scan(&id, &name, &email, &phone, &dob, &salary)
		results = append(results, Employee{
			Name:   name,
			Email:  email,
			Phone:  phone,
			DOB:    dob,
			Salary: salary,
		})
	}
	rows.Close()

	return results
}

func main() {
	openLogFile("/home/whitehat/access.log")

	tmpl := template.Must(template.ParseFiles("index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		data := EmployeePageData{
			PageTitle: "Employees",
			Employees: getEmployees(req.PostFormValue("search")),
		}
		tmpl.Execute(w, data)
	})

	http.ListenAndServe(":"+os.Getenv("VIRTUAL_PORT"), logHandler(http.DefaultServeMux))
}

func logHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		x, err := httputil.DumpRequest(r, true)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
		log.Println(fmt.Sprintf("Q %q", x))
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, r)
		log.Println(fmt.Sprintf("A %d %q", rec.Code, rec.Body))

		handler.ServeHTTP(w, r)
	})
}

func openLogFile(logfile string) {
	lf, _ := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)
	log.SetOutput(lf)
	log.SetFlags(log.Ltime)
}
