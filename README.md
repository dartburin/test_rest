# test_rest
Simple RESTful API application for work with PostgreSQL DB.

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