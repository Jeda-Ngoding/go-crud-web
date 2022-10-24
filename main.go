package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Product struct {
	Id          int
	Name        string
	Description string
	Price       float64
	Created_At  time.Time
	Updated_At  time.Time
}

type PageData struct {
	PageTitle string
	Products  []Product
}

func connectDB() *sql.DB {
	db, err := sql.Open("mysql", "root:password@(127.0.0.1:3306)/go_crud_web?parseTime=true")

	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db

}

func index(w http.ResponseWriter, r *http.Request) {

	rows, err := connectDB().Query(`SELECT id, name, description,price FROM products`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	data := PageData{
		PageTitle: "Products",
		Products:  []Product{},
	}

	for rows.Next() {
		var p Product

		err := rows.Scan(&p.Id, &p.Name, &p.Description, &p.Price)
		if err != nil {
			log.Fatal(err)
		}
		data.Products = append(data.Products, p)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	var view = template.Must(template.ParseFiles("views/index.html"))
	view.Execute(w, data)
}

func create(w http.ResponseWriter, r *http.Request) {
	var view = template.Must(template.ParseFiles("views/create.html"))
	view.Execute(w, nil)
}

func edit(w http.ResponseWriter, r *http.Request) {
	var view = template.Must(template.ParseFiles("views/edit.html"))
	view.Execute(w, nil)
}

func store(w http.ResponseWriter, r *http.Request) {
	var view = template.Must(template.ParseFiles("views/edit.html"))
	view.Execute(w, nil)
}

func update(w http.ResponseWriter, r *http.Request) {
	var view = template.Must(template.ParseFiles("views/edit.html"))
	view.Execute(w, nil)
}

func delete(w http.ResponseWriter, r *http.Request) {
	var view = template.Must(template.ParseFiles("views/edit.html"))
	view.Execute(w, nil)
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("Error Load Environment")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/create", create)
	mux.HandleFunc("/store", store)
	mux.HandleFunc("/edit", edit)
	mux.HandleFunc("/update", update)
	mux.HandleFunc("/delete", delete)

	http.ListenAndServe(":"+port, mux)

}
