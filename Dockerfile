FROM golang:1.15
ADD . /go/src/github.com/cogolabs/beyond
RUN go get -x github.com/cogolabs/beyond/cmd/httpd
CMD ["httpd", "--help"]
