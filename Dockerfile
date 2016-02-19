FROM golang:1.5.3

ADD . /go/src/github.com/aubm/cats-api

RUN go get github.com/gorilla/mux
RUN go get gopkg.in/mgo.v2
RUN go install github.com/aubm/cats-api

ENTRYPOINT /go/bin/cats-api

EXPOSE 8080
