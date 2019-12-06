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

// Data for configuration of DB connecttion
type Config struct {
	User string
	Pass string
	Db   string
	Host string
	Port string
}

// Book info
type Book struct {
	Id     int64
	Author string
	Title  string
}

// Create connection to db
func ConnectToDB(conf Config) (ParamDB, error) {
	var err error = nil
	var par ParamDB
	par.Conf = conf
	par.Base = nil

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
	queue := fmt.Sprintf("SELECT COUNT(*) FROM pg_database WHERE datname = '%s';", name)

	row := base.QueryRow(queue)

	var cnt int
	if err := row.Scan(&cnt); err != nil {
		panic(err)
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
		"title varchar(50)," +
		"author varchar(50)" +
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

// Insert book info into database
func (b *Book) InsertBook(base *sql.DB) (int64, error) {
	if base == nil {
		return 0, errors.New("DB not opened")
	}

	query := fmt.Sprintf("INSERT INTO BookInfo (title, author) VALUES('%s', '%s') RETURNING id;", b.Title, b.Author)
	row := base.QueryRow(query)

	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, errors.New("Error insert book - bad id")
	}
	return id, nil
}

// Insert book info into database
func (b *Book) DeleteBook(base *sql.DB) error {
	if base == nil {
		return errors.New("DB not opened")
	}

	if b.Id <= 0 {
		return errors.New("Id not set")
	}

	query := fmt.Sprintf("DELETE FROM BookInfo where id = %v;", b.Id)
	_, err := base.Exec(query)

	if err != nil {
		return errors.New("Error delete book")
	}
	return nil
}

// Update book info into database
func (b *Book) UpdateBook(base *sql.DB) error {
	if base == nil {
		return errors.New("DB not opened")
	}

	if b.Id <= 0 {
		return errors.New("Id not set")
	}

	var query string
	switch {
	case b.Author != "" && b.Title != "":
		query = fmt.Sprintf("UPDATE BookInfo SET author = '%s', title = '%s' WHERE id = %v;", b.Author, b.Title, b.Id)
	case b.Author != "":
		query = fmt.Sprintf("UPDATE BookInfo SET author = '%s' WHERE id = %v;", b.Author, b.Id)
	case b.Title != "":
		query = fmt.Sprintf("UPDATE BookInfo SET title = '%s' WHERE id = %v;", b.Title, b.Id)
	default:
		return nil
	}
	fmt.Println(query)
	_, err := base.Exec(query)

	if err != nil {
		return errors.New("Error update book")
	}
	return nil
}

// Select book info from database
func (b *Book) SelectBook(base *sql.DB) (*[]Book, error) {
	bb := make([]Book, 0)
	if base == nil {
		return &bb, errors.New("DB not opened")
	}

	var query string
	if b.Id > 0 {
		query = fmt.Sprintf("SELECT id, title, author FROM BookInfo WHERE id = %v ORDER BY id;", b.Id)
	} else {
		query = fmt.Sprintf("SELECT id, title, author FROM BookInfo ORDER BY id;")
	}

	rows, err := base.Query(query)
	if err != nil {
		return &bb, errors.New("Select error")
	}
	defer rows.Close()

	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.Id, &book.Title, &book.Author); err != nil {
			return &bb, errors.New("Select scan error")
		}
		bb = append(bb, book)
	}

	if err != nil {
		return &bb, errors.New("Error update book")
	}
	return &bb, nil
}
