package servgrpc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"

	bk "test_rest/cmd/grpc/books"
	pdb "test_rest/cmd/rest/postgredb"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

// Data for handlers
type supportGRPC struct {
	base  *sql.DB
	Host  string
	Port  string
	GPort string
	Conn  string
	GConn string
	wg    sync.WaitGroup
}

// New creates new server greeter
func New(b *sql.DB, host string, port string, gport string) *supportGRPC {
	str := fmt.Sprintf("%s:%s", host, port)
	gstr := fmt.Sprintf("%s:%s", host, gport)
	return &supportGRPC{
		base:  b,
		Host:  host,
		Port:  port,
		GPort: gport,
		Conn:  str,
		GConn: gstr,
	}
}

// Start server gRPC + gate REST
func (s *supportGRPC) Start() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		log.Fatal(s.startGRPC())
	}()

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		log.Fatal(s.startREST())
	}()

	s.wg.Wait()
}

func (s *supportGRPC) startGRPC() error {
	lis, err := net.Listen("tcp", s.GConn)
	if err != nil {
		return err
	}
	srv := grpc.NewServer()
	bk.RegisterLibraryServer(srv, s)

	err = srv.Serve(lis)
	if err != nil {
		return err
	}
	return nil
}

func (s *supportGRPC) startREST() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := bk.RegisterLibraryHandlerFromEndpoint(ctx, mux, s.GConn, opts)
	if err != nil {
		return err
	}

	return http.ListenAndServe(s.Conn, mux)
}

func (s *supportGRPC) GetBooks(ctx context.Context, imsg *bk.GetBookRequest) (*bk.Books, error) {
	var bb pdb.Book
	books, err := bb.SelectBook(s.base)
	if err != nil {
		return nil, errors.New("Error get book")
	}

	bs := bk.Books{}
	bs.Books = make([]*bk.OneBook, 0, 50)
	for i := range *books {
		b := bk.OneBook{}
		b.Id = (*books)[i].Id
		b.Author = (*books)[i].Author
		b.Title = (*books)[i].Title
		bs.Books = append(bs.Books, &b)
	}

	return &bs, nil
}

// Handler for post request
func (s *supportGRPC) PostBook(ctx context.Context, imsg *bk.PostBookRequest) (*bk.OneBook, error) {

	b := &pdb.Book{}
	b.Author = imsg.Message.Author
	b.Title = imsg.Message.Title

	//fmt.Printf("Add %v\n", b)

	id, err := b.InsertBook(s.base)
	if err != nil {
		return nil, errors.New("Error post (insert) book")
	}

	b.Id = id
	// Maybe select not need
	_, err = b.SelectBook(s.base)
	if err != nil {
		return nil, errors.New("Error post (select) book")
	}

	bb := bk.OneBook{}
	bb.Id = id
	bb.Author = b.Author
	bb.Title = b.Title
	//fmt.Printf("Add ret %v\n", bb)

	return &bb, nil
}

// Handler for delete request
func (s *supportGRPC) DeleteBook(ctx context.Context, imsg *bk.DeleteBookRequest) (*bk.Result, error) {
	id, err := strconv.Atoi(imsg.BookId)
	if err != nil {
		return nil, errors.New("Error delete book (bad id)")
	}
	b := &pdb.Book{}
	b.Id = int64(id)

	err = b.DeleteBook(s.base)
	if err != nil {
		return nil, errors.New("Error delete book")
	}

	bb := bk.Result{}
	bb.Rez = fmt.Sprintf("Delete book with id = %v", b.Id)
	return &bb, nil
}

// Handler for put request
func (s *supportGRPC) UpdateBook(ctx context.Context, imsg *bk.UpdateBookRequest) (*bk.OneBook, error) {
	id, err := strconv.Atoi(imsg.BookId)
	if err != nil {
		return nil, errors.New("Error delete book (bad id)")
	}
	b := &pdb.Book{}
	b.Id = int64(id)
	b.Author = imsg.Message.Author
	b.Title = imsg.Message.Title

	if b.Author == "" || b.Title == "" {
		return nil, errors.New("Error some parameters not set for PUT request")
	}

	err = b.UpdateBook(s.base)
	if err != nil {
		return nil, errors.New("Error update book")
	}

	bb := bk.OneBook{}
	bb.Id = b.Id
	bb.Author = b.Author
	bb.Title = b.Title

	return &bb, nil
}

// Handler for patch request
func (s *supportGRPC) PathBook(ctx context.Context, imsg *bk.UpdateBookRequest) (*bk.OneBook, error) {
	id, err := strconv.Atoi(imsg.BookId)
	if err != nil {
		return nil, errors.New("Error delete book (bad id)")
	}
	b := &pdb.Book{}
	b.Id = int64(id)
	b.Author = imsg.Message.Author
	b.Title = imsg.Message.Title

	err = b.UpdateBook(s.base)
	if err != nil {
		return nil, errors.New("Error update book")
	}

	bb := bk.OneBook{}
	bb.Id = b.Id
	bb.Author = b.Author
	bb.Title = b.Title

	return &bb, nil
}
