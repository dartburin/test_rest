FROM golang:latest

ARG A_DB_USER=postgres
ARG A_DB_PASS=postgres
ARG A_DB_BASE=books
ARG A_DB_HOST=servdb
ARG A_DB_PORT=5432
ARG A_HTTP_HOST=0.0.0.0
ARG A_HTTP_PORT=8080


ENV APP_NAME rest


ENV DB_USER=${A_DB_USER} 
ENV DB_PASS=${A_DB_PASS} 
ENV DB_BASE=${A_DB_BASE} 
ENV DB_HOST=${A_DB_HOST} 
ENV DB_PORT=${A_DB_PORT} 
ENV HTTP_HOST=${A_HTTP_HOST} 
ENV HTTP_PORT=${A_HTTP_PORT} 


#RUN go version

RUN mkdir -p ${GOPATH}/src/${APP_NAME}

WORKDIR /go/src/${APP_NAME}

COPY go.mod ./
COPY go.sum ./

COPY ./cmd ./
COPY ./cmd/grpc/ ./cmd/
COPY ./cmd/grpc/* ./cmd/grpc/
COPY ./cmd/grpc/servgrpc/* ./cmd/grpc/servgrpc/
COPY ./cmd/grpc/books/* ./cmd/grpc/books/

COPY ./cmd/rest/ ./cmd/
COPY ./cmd/rest/* ./cmd/rest/
COPY ./cmd/rest/postgredb/* ./cmd/rest/postgredb/
COPY ./cmd/rest/servhttp/* ./cmd/rest/servhttp/

#RUN ls -l
#RUN pwd

RUN go mod download 
#RUN protoc -I/usr/local/include -I. -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:. cmd/grpc/books/books.proto
#RUN protoc -I/usr/local/include -I. -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:. cmd/grpc/books/books.proto
RUN go build cmd/grpc/main.go
RUN ls -l

EXPOSE ${HTTP_PORT}
EXPOSE ${DB_PORT}


CMD ./main -dbhost=${DB_HOST} -dbbase=${DB_BASE} -dbuser=${DB_USER} -dbpass=${DB_PASS} -dbport=${DB_PORT} -httphost=${HTTP_HOST} -httpport=${HTTP_PORT}

