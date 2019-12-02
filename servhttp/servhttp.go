package servhttp

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	pdb "test_rest/postgredb"
)

// Intenface for not numered requests
type ParamsServer struct {
	base *sql.DB
	url  string
}

// Intenface for numered requests
type ParamsServerNum struct {
	base *sql.DB
	url  string
}

// Method create a server
func HttpServer(b *sql.DB) {
	params := &ParamsServer{base: b, url: "/books"}
	http.Handle("/books", params)

	paramsNum := &ParamsServerNum{base: b, url: "/books/"}
	http.Handle("/books/", paramsNum)

	fmt.Printf("Starting HTTP server at port 8080\n")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("err=%v\n", err)
	}
}

// Handler for not numered requests
func (p *ParamsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
		var bb pdb.Book
		books, err := bb.SelectBook(p.base)
		if err != nil {
			fmt.Fprintf(w, "Database error")
		} else {
			res, err := json.Marshal(books)
			if err != nil {
				fmt.Fprintf(w, "JSON parse error = %v!\n", err)
			} else {
				fmt.Fprintf(w, "%s", string(res))
			}
		}

	case r.Method == http.MethodPost:
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(w, "POST reading error\n")
		} else {
			b := &pdb.Book{}
			err = json.Unmarshal(data, b)
			if err != nil {
				fmt.Fprintf(w, "POST JSON parsing error\n")
			} else {
				id, err := b.InsertBook(p.base)
				if err != nil {
					fmt.Fprintf(w, "Database insert error\n")
				} else {
					b.Id = id
					// Maybe select not need
					book, err := b.SelectBook(p.base)
					if err != nil {
						fmt.Fprintf(w, "Database reading new record error\n")
					} else {
						res, err := json.Marshal(book)
						if err != nil {
							fmt.Fprintf(w, "JSON parse error = %v !\n", err)
						} else {
							fmt.Fprintf(w, "%s", string(res))
						}
					}
				}
			}

		}

	default:
		fmt.Fprintf(w, "Method not used\n")
	}
}

// Handler for numered requests
func (p *ParamsServerNum) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	str := strings.Replace(r.URL.String(), p.url, "", 1)
	id, err := strconv.Atoi(str)
	if err != nil {
		fmt.Fprintf(w, "Not valid Id in request (%s)\n", str)
		return
	}

	switch {
	case r.Method == http.MethodDelete:
		b := &pdb.Book{}
		b.Id = int64(id)
		err = b.DeleteBook(p.base)
		if err != nil {
			fmt.Fprintf(w, "Database delete error\n")
		} else {
			fmt.Fprintf(w, "Book %v deleted from database\n", id)
		}

	case r.Method == http.MethodPut:
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(w, "PATCH reading error\n")
		} else {
			b := &pdb.Book{}
			err = json.Unmarshal(data, b)
			if err != nil {
				fmt.Fprintf(w, "PATCH JSON parsing error\n")
			} else {
				b.Id = int64(id)
				if b.Author == "" || b.Title == "" {
					fmt.Fprintf(w, "Some parameters not set for PUT request.\n")
				} else {
					err = b.UpdateBook(p.base)
					if err != nil {
						fmt.Fprintf(w, "Database update error.\n")
					} else {
						fmt.Fprintf(w, "Update database record id=%v\n", id)
					}

				}
			}

		}

	case r.Method == http.MethodPatch:
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(w, "PATCH reading error\n")
		} else {
			b := &pdb.Book{}
			err = json.Unmarshal(data, b)
			if err != nil {
				fmt.Fprintf(w, "PATCH JSON parsing error\n")
			} else {
				b.Id = int64(id)
				err = b.UpdateBook(p.base)
				if err != nil {
					fmt.Fprintf(w, "Database update error\n")
				} else {
					fmt.Fprintf(w, "Update database record id=%v\n", id)
				}
			}

		}

	default:
		fmt.Fprintf(w, "Method not used\n")
	}
}
