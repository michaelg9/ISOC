FROM golang:1.6

RUN mkdir -p /go/src/github.com/michaelg9/ISOC/server
WORKDIR /go/src/github.com/michaelg9/ISOC/server

CMD ["go-wrapper", "run"]

COPY . /go/src/github.com/michaelg9/ISOC/server
RUN go get github.com/codegangsta/gin
RUN go-wrapper download
RUN go-wrapper install

EXPOSE 3000
