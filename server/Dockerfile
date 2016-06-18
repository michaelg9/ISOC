FROM golang:1.6

RUN mkdir -p /go/src/github.com/michaelg9/ISOC/server
WORKDIR /go/src/github.com/michaelg9/ISOC/server

CMD ["go-wrapper", "run"]

COPY . /go/src/github.com/michaelg9/ISOC/server
RUN go get github.com/codegangsta/gin
RUN go-wrapper download
RUN go-wrapper install

# Set the port on which the server will run
ENV PORT 8080

EXPOSE 3000
