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

COPY ./rest.go ./
COPY ./postgredb ./postgredb/
COPY ./servhttp ./servhttp/

#RUN ls -l
#RUN pwd

#RUN go get -d -v
RUN go mod download 
RUN go build ${APP_NAME}.go
RUN ls -l

EXPOSE ${HTTP_PORT}
EXPOSE ${DB_PORT}


CMD ./${APP_NAME} -dbhost=${DB_HOST} -dbbase=${DB_BASE} -dbuser=${DB_USER} -dbpass=${DB_PASS} -dbport=${DB_PORT} -httphost=${HTTP_HOST} -httpport=${HTTP_PORT}

