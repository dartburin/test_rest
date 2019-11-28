package main

import (
	"fmt"
	pdb "test_rest/postgredb"
)

func main() {
	fmt.Println("RESTful service !!!")

	config := pdb.Config{
		User: "postgres",
		Pass: "postgres",
		Db:   "books",
		Host: "127.0.0.1",
		Port: "5432",
	}

	param, err := pdb.ConnectToDB(config)
	if err != nil {
		return
	}

	fmt.Printf("Connect ! %v \n\n", param)
	param.Base.Close()

	fmt.Println("RESTful service end !!!\n")
}
