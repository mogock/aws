package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" //Without this you receive := unknown driver "postgres" (forgotten import?) , pq.init() := for Alias
	"log"
	"net/http"
	"strconv"
	resp "mogock.com/bookstore/response"
)

type Book struct {
	isbn string
	title string
	author string
	price float32
}

var db *sql.DB

//Stard Database Connection
func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://postgres:123456@localhost/postgres?sslmode=disable")	//Se usa igual si las variable fueron creadas
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	if err = db.Ping(); err != nil { //Muy Interensante
		log.Fatal(err)
	}
}

///books
func booksIndex (w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	rows, err := db.Query("SELECT * FROM books")
	if err != nil {
		http.Error(w, http.StatusText(500),  500)
		return
	}
	defer rows.Close()

	bks := make([]*Book, 0)
	for rows.Next() {
		bk := new(Book)
		err := rows.Scan(&bk.isbn, &bk.title, &bk.author, &bk.price)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		bks = append(bks, bk)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	//If Everything is ok? then :=?
	for _, bk := range bks {
		fmt.Fprintf(w,"%s, %s, %s, £%.2f\n", bk.isbn, bk.title, bk.author, bk.price)
	}

}

// --  /books/show?isbn=978-1505255607
func booksShow(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	isbn := r.FormValue("isbn")
	if isbn == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	row := db.QueryRow("SELECT * FROM books WHERE isbn = $1", isbn) //This is a single row query
	bk := new(Book)
	err := row.Scan(&bk.isbn, &bk.title, &bk.author, &bk.price)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}else if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	fmt.Fprintf(w, "%s, %s, %s, £%.2f\n", bk.isbn, bk.title, bk.author, bk.price)
}

//--- /books/create
func booksCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 500)
		return
	}
	isbn := r.FormValue("isbn")
	title := r.FormValue("title")
	author := r.FormValue("author")
	if isbn == "" || title == "" || author == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	price, err := strconv.ParseFloat(r.FormValue("price"), 32)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	result, err := db.Exec("INSERT INTO books VALUES($1, $2, $3, $4)", isbn, title, author, price)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	fmt.Fprintf(w, "Book %s created successfully (%d row affected)\n", isbn, rowsAffected)
}

func main() {
	//http.HandleFunc("/books", booksIndex)
	//http.HandleFunc("/books/show", booksShow)
	//http.HandleFunc("/books/create", booksCreate)
	//http.ListenAndServe(":8080", nil)
	fmt.Println(resp.GLOBAL)
	http.HandleFunc("/json", resp.FooJSON)
	http.HandleFunc("/xml", resp.FooXML)
	http.ListenAndServe(":8080", nil)
}