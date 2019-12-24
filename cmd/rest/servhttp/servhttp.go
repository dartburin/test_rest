package servhttp

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	pdb "test_rest/cmd/rest/postgredb"

	"github.com/gorilla/mux"
)

// Data for configuration DB connecttion
type ConnectConfig struct {
	Host string
	Port string
}

// Data for handlers
type supportHTTP struct {
	base *sql.DB
}

// Handler for get request
func (d *supportHTTP) GetBooks() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var bb pdb.Book
		books, err := bb.SelectBook(d.base)
		if err != nil {
			w.WriteHeader(500) // Internal Server Error
			return
		}

		res, err := json.Marshal(books)
		if err != nil {
			w.WriteHeader(500) // Internal Server Error
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200) // OK
		fmt.Fprintf(w, "%s", string(res))
	}
}

// Handler for post request
func (d *supportHTTP) PostBook() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(500) // Internal Server Error
			return
		}

		b := &pdb.Book{}
		err = json.Unmarshal(data, b)
		if err != nil {
			w.WriteHeader(400) // Bad Request
			return
		}

		id, err := b.InsertBook(d.base)
		if err != nil {
			w.WriteHeader(500) // Internal Server Error
			return
		}

		b.Id = id
		// Maybe select not need
		book, err := b.SelectBook(d.base)
		if err != nil {
			w.WriteHeader(400) // Bad Request
			return
		}

		res, err := json.Marshal(book)
		if err != nil {
			w.WriteHeader(500) // Internal Server Error
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200) // OK
		fmt.Fprintf(w, "%s", string(res))
	}
}

// Handler for delete request
func (d *supportHTTP) DelBook() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			w.WriteHeader(400) // Bad Request
			return
		}
		b := &pdb.Book{}
		b.Id = int64(id)

		err = b.DeleteBook(d.base)
		if err != nil {
			w.WriteHeader(500) // Internal Server Error
			return
		}

		w.WriteHeader(200) // OK
		//fmt.Fprintf(w, "Book %v deleted from database.\n", id)
	}
}

// Handler for put request
func (d *supportHTTP) PutBook() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(400) // Bad Request
			return
		}

		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			w.WriteHeader(400) // Bad Request
			return
		}

		b := &pdb.Book{}
		err = json.Unmarshal(data, b)
		if err != nil {
			w.WriteHeader(400) // Bad Request
			return
		}

		b.Id = int64(id)
		if b.Author == "" || b.Title == "" {
			//fmt.Fprintf(w, "Some parameters not set for PUT request.\n")
			w.WriteHeader(400) // Bad Request
			return
		}

		err = b.UpdateBook(d.base)
		if err != nil {
			w.WriteHeader(500) // Internal Server Error
			return
		}

		w.WriteHeader(200) // OK
		//fmt.Fprintf(w, "Update database record id=%v\n", id)
	}
}

// Handler for patch request
func (d *supportHTTP) PatchBook() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(400) // Bad Request
			return
		}

		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			w.WriteHeader(400) // Bad Request
			return
		}
		b := &pdb.Book{}
		b.Id = int64(id)

		err = json.Unmarshal(data, b)
		if err != nil {
			w.WriteHeader(400) // Bad Request
			return
		}

		b.Id = int64(id)
		err = b.UpdateBook(d.base)
		if err != nil {
			w.WriteHeader(500) // Internal Server Error
			return
		}

		w.WriteHeader(200) // OK
		//fmt.Fprintf(w, "Update database record id=%v\n", id)
	}
}

// Method create a server
func (c *ConnectConfig) HttpServer(b *sql.DB) {
	sup := &supportHTTP{base: b}

	router := mux.NewRouter()
	router.HandleFunc("/books", sup.GetBooks()).Methods("GET")
	router.HandleFunc("/books", sup.PostBook()).Methods("POST")
	router.HandleFunc("/books/{id}", sup.DelBook()).Methods("DELETE")
	router.HandleFunc("/books/{id}", sup.PutBook()).Methods("PUT")
	router.HandleFunc("/books/{id}", sup.PatchBook()).Methods("PATCH")

	fmt.Printf("Starting HTTP server at %s:%s\n", c.Host, c.Port)

	err := http.ListenAndServe(c.Host+":"+c.Port, router)
	if err != nil {
		fmt.Printf("err=%v\n", err)
	}
}
