package main

import (
	"database/sql"
	"fmt"
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
	Image       string
	Status      int
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

	rows, err := connectDB().Query(`SELECT id, name, description,image,price,status FROM products`)
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

		err := rows.Scan(&p.Id, &p.Name, &p.Description, &p.Image, &p.Price, &p.Status)
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

	data := PageData{
		PageTitle: "Create Product",
	}

	var view = template.Must(template.ParseFiles("views/create.html"))
	view.Execute(w, data)
}

func edit(w http.ResponseWriter, r *http.Request) {
	var view = template.Must(template.ParseFiles("views/edit.html"))
	view.Execute(w, nil)
}

func store(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/create", r.Response.StatusCode)
	}

	product := Product{
		Id:          0,
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Price:       0,
		Image:       "-",
		Status:      1,
		Created_At:  time.Now(),
		Updated_At:  time.Now(),
	}

	result, err := connectDB().Exec(`INSERT INTO products (id,name, description,price,image,status, created_at,updated_at) VALUES (?,?, ?, ?,?, ?, ?,?)`, nil, product.Name, product.Description, product.Price, product.Image, product.Status, product.Created_At, product.Updated_At)
	if err != nil {
		log.Fatal(err)
	}

	id, err := result.LastInsertId()

	if id != 0 {
		http.Redirect(w, r, "/", 200)
		return
	} else {
		fmt.Println(err)
		http.Redirect(w, r, "/create", r.Response.StatusCode)
		return
	}

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
