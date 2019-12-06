package main

import (
	"flag"
	"fmt"
	"os"
	pdb "test_rest/postgredb"
	rest "test_rest/servhttp"
)

func main() {
	// Load init parameters
	var configDB pdb.Config
	dbHostName := flag.String("dbhost", "", "host name")
	dbUser := flag.String("dbuser", "", "user db name")
	dbPass := flag.String("dbpass", "", "user db pass")
	dbBase := flag.String("dbbase", "", "database name")
	dbPort := flag.String("dbport", "", "port for batabase connect")

	var configRest rest.ConnectConfig
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
	configRest.Port = *httpPort
	configRest.Host = *httpHostName

	configDB.Port = *dbPort
	configDB.Host = *dbHostName
	configDB.User = *dbUser
	configDB.Pass = *dbPass
	configDB.Db = *dbBase

	// Connect to DB
	parDB, err := pdb.ConnectToDB(configDB)
	if err != nil {
		return
	}
	defer parDB.Base.Close()

	//Start HTTP server
	configRest.HttpServer(parDB.Base)
}
