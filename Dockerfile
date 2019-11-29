FROM golang:latest

ENV APP_NAME rest
ENV PORT 8080
ENV PORT2 5432
ENV POSTGRE_HOST servdb

RUN go version

RUN mkdir -p ${GOPATH}/src/${APP_NAME}

WORKDIR /go/src/${APP_NAME}

COPY go.mod ./
COPY go.sum ./

COPY ./rest.go ./
COPY ./postgredb ./postgredb/

RUN ls -l
RUN pwd
RUN echo ${POSTGRE_HOST}

#RUN go get -d -v
RUN go mod download 
RUN go build ${APP_NAME}.go
RUN ls -l

CMD ./${APP_NAME} -host=${POSTGRE_HOST}

EXPOSE ${PORT}
EXPOSE ${PORT2}
