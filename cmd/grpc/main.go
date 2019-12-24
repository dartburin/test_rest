package main

import (
	"flag"
	"fmt"
	"os"
	gp "test_rest/cmd/grpc/servgrpc"
	pdb "test_rest/cmd/rest/postgredb"
)

func main() {
	//fmt.Printf("Server init \n")
	// Load init parameters
	var configDB pdb.Config
	dbHostName := flag.String("dbhost", "", "host name")
	dbUser := flag.String("dbuser", "", "user db name")
	dbPass := flag.String("dbpass", "", "user db pass")
	dbBase := flag.String("dbbase", "", "database name")
	dbPort := flag.String("dbport", "", "port for batabase connect")

	httpHostName := flag.String("httphost", "", "host name")
	httpPort := flag.String("httpport", "", "port for http connect")

	flag.Parse()

	// Check existing obligatory http and db parameters
	if *httpHostName == "" || *httpPort == "" ||
		*dbHostName == "" || *dbUser == "" || *dbPass == "" ||
		*dbBase == "" || *dbPort == "" {
		fmt.Println("Init error: set not all obligatory parameters.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Set http and db parameters
	configDB.Port = *dbPort
	configDB.Host = *dbHostName
	configDB.User = *dbUser
	configDB.Pass = *dbPass
	configDB.Db = *dbBase

	// Connect to DB
	//fmt.Printf("Base start\n")
	parDB, err := pdb.ConnectToDB(configDB)
	if err != nil {
		fmt.Printf("Base not connect\n")
		os.Exit(1)
	}
	defer parDB.Base.Close()

	//Start gRPC + HTTP server
	g := gp.New(parDB.Base, *httpHostName, *httpPort, "8086")
	fmt.Printf("Server start\n")
	g.Start()
}
