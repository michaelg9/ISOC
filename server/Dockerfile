FROM golang:1.6

# Setup the filesystem
RUN mkdir -p /go/src/github.com/michaelg9/ISOC/server
WORKDIR /go/src/github.com/michaelg9/ISOC/server

CMD ["go-wrapper", "run"]

COPY . /go/src/github.com/michaelg9/ISOC/server

# Install utility functions for development
# gin for automatic recompiling and restarting the server, assert for testing
RUN go get github.com/codegangsta/gin
RUN go get github.com/stretchr/testify/assert

# Install the application
RUN go-wrapper download
RUN go-wrapper install

EXPOSE 3000
