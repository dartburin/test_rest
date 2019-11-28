FROM golang:latest

ENV APP_NAME rest
ENV PORT 8080
ENV PORT2 5432

RUN go version

RUN mkdir -p ${GOPATH}/src/${APP_NAME}

WORKDIR /go/src/${APP_NAME}

COPY go.mod ./
COPY go.sum ./

COPY ./rest.go ./
COPY ./postgredb ./postgredb/

RUN ls -l
RUN pwd

#RUN go get -d -v
RUN go mod download 
RUN go build ${APP_NAME}.go
RUN ls -l

CMD ./${APP_NAME}

EXPOSE ${PORT}
EXPOSE ${PORT2}
