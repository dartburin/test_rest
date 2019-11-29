package main

import (
	"flag"
	"fmt"
	pdb "test_rest/postgredb"
)

func main() {
	fmt.Println("RESTful service !!!")

	config := pdb.Config{
		User: "postgres",
		Pass: "postgres",
		Db:   "books",
		Host: "localhost",
		//		Host: "127.0.0.1",
		//		Host: "servdb",
		Port: "5432",
	}

	hostName := flag.String("host", "", "a host name")
	flag.Parse()

	// Check existing host name param
	if *hostName != "" {
		config.Host = *hostName
		fmt.Printf("Set postgreSQL server host = %v \n", *hostName)
	}

	parDB, err := pdb.ConnectToDB(config)
	if err != nil {
		return
	}
	defer parDB.Base.Close()

	fmt.Printf("Connect ! %v \n\n", parDB)

	//insert test
	bb := pdb.Book{
		Id:     0,
		Title:  "Best book 1",
		Author: "Good author 1",
	}
	id, err := bb.InsertBook(parDB.Base)
	fmt.Printf("Book inserted %v (err = %v) !!!\n", id, err)

	bb.Title = "Best book 2"
	bb.Author = "Good author 2"
	id, err = bb.InsertBook(parDB.Base)
	fmt.Printf("Book inserted %v (err = %v) !!!\n", id, err)
	toDel := id

	bb.Title = "Best book 3"
	bb.Author = "Good author 3"
	id, err = bb.InsertBook(parDB.Base)
	fmt.Printf("Book inserted %v (err = %v) !!!\n", id, err)
	toSel := id

	bb.Title = "Best book 4"
	bb.Author = "Good author 4"
	id, err = bb.InsertBook(parDB.Base)
	fmt.Printf("Book inserted %v (err = %v) !!!\n", id, err)
	toUpd := id

	//delete test
	bb.Id = toDel
	err = bb.DeleteBook(parDB.Base)
	fmt.Printf("Book deleted (err = %v) !!!\n", err)

	//update test
	bb.Id = toUpd
	bb.Author = "Very best author"
	bb.Title = "111"
	err = bb.UpdateBook(parDB.Base)
	fmt.Printf("Book updated (err = %v) !!!\n\n", err)

	//select one test
	bb.Id = toSel
	books, err := bb.SelectBook(parDB.Base)
	fmt.Printf("Book one select \n%v (err = %v) !!!\n\n", books, err)

	//select lot test
	bb.Id = 0
	books, err = bb.SelectBook(parDB.Base)
	fmt.Printf("Books lot select \n%v (err = %v) !!!\n\n", books, err)

	fmt.Printf("RESTful service end !!!\n")
}
