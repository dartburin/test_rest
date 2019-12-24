# test_rest
Simple RESTful API application for work with PostgreSQL DB.

(Local start)
server:  go build cmd/grpc/main.go
./main   -dbbase=books  -dbhost=servdb  -dbpass=postgres  -dbport=5432  -dbuser=postgres   -httphost=0.0.0.0  -httpport=8080


(START)
docker-compose up --build


Sample of testing orders
------------------------

(GET)
http://localhost:8080/books

(POST)
curl -XPOST "http://localhost:8080/books" -d '{"Author": "Some author", "Title": "New book 1"}'

(DELETE)
curl -XDELETE "http://localhost:8080/books/91" 

(PATCH)
curl -XPATCH "http://localhost:8080/books/43" -d '{"Author": "A100", "Title": "Book 1"}'

(PUT)
curl -XPUT "http://localhost:8080/books/43" -d '{"Author": "A1", "Title": "Book 1"}'