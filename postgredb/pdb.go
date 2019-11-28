package postgredb

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

// DataBase info
type ParamDB struct {
	Conf Config
	Base *sql.DB
}

// Data for configuration DB connecttion
type Config struct {
	User string
	Pass string
	Db   string
	Host string
	Port string
}

// Creae connection to db
func ConnectToDB(conf Config) (ParamDB, error) {
	var par ParamDB
	par.Conf = conf
	par.Base = nil
	err := errors.New("No error.")

	if conf.User == "" || conf.Pass == "" || conf.Db == "" ||
		conf.Host == "" || conf.Port == "" {
		err = errors.New("Bad connection parameters")
		return par, err
	}

	// Check connect to default server PostgreSQL database
	connstr := fmt.Sprintf("user=%s password=%s host=%s port=%s sslmode=disable",
		conf.User, conf.Pass, conf.Host, conf.Port)
	defdb, err := openDB(connstr)
	if err != nil {
		err = errors.New("Bad connection to server db")
		return par, err
	}
	defer defdb.Close()

	// Check user database exists and create if need
	if err = checkExistsUserDB(defdb, conf.Db); err != nil {
		par, err = createDB(conf, defdb)

		if err != nil {
			fmt.Printf("Error: %v !\n", err.Error())
			defer par.Base.Close()
			par.Base = nil
			return par, err
		}

		err = createBookTable(par.Base)
		if err != nil {
			fmt.Printf("Error: %v !\n", err.Error())
			defer par.Base.Close()
			par.Base = nil
			return par, err
		}
	} else {
		// Connect to user database
		connstr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
			conf.User, conf.Pass, conf.Host, conf.Port, conf.Db)
		db, err := openDB(connstr)
		if err != nil {
			err = errors.New("Bad connection to user db")
			return par, err
		}
		par.Base = db
	}

	return par, nil
}

// Close user database
func (par *ParamDB) Close() (err error) {
	if par.Base == nil {
		return nil
	}
	fmt.Printf("Close ! %v \n\n", par)

	if err = par.Base.Close(); err != nil {
		return err
	}

	par.Base = nil
	return nil
}

// Check user database exists and create if need
func checkExistsUserDB(base *sql.DB, name string) error {
	queue := fmt.Sprintf("select count(*) from pg_database where datname = '%s';", name)

	rows, err := base.Query(queue)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var cnt int
	for rows.Next() {
		if err := rows.Scan(&cnt); err != nil {
			panic(err)
		}
		break
	}
	if cnt != 1 {
		return errors.New("No base")
	}
	return nil
}

// Create user database
func createDB(conf Config, base *sql.DB) (ParamDB, error) {
	var par ParamDB
	par.Conf = conf
	par.Base = nil

	queue := fmt.Sprintf("CREATE DATABASE %s;", conf.Db)

	_, err := base.Exec(queue)
	if err != nil {
		panic(err)
	}

	// Check connect to user DB
	connstr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		conf.User, conf.Pass, conf.Host, conf.Port, conf.Db)

	db, err := openDB(connstr)
	if err != nil {
		err = errors.New("Bad connection to user db")
		return par, err
	}

	par.Base = db
	return par, nil
}

// Create BookInfo table
func createBookTable(base *sql.DB) error {
	_, err := base.Exec("CREATE TABLE IF NOT EXISTS BookInfo(" +
		"id SERIAL PRIMARY KEY," +
		"name varchar(50)," +
		"autor varchar(50)" +
		");")

	if err != nil {
		panic(err)
	}

	return nil
}

// Check connect to DB
func openDB(confstr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", confstr)

	if err != nil {
		err = errors.New("Bad connection to db")
		return nil, err
	}

	// Ping the database
	if err = db.Ping(); err != nil {
		err = errors.New("Couldn't ping to database")
		db.Close()
		return nil, err
	}
	return db, nil
}
