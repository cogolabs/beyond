FROM golang:1.11
ADD . /go/src/github.com/cogolabs/beyond
RUN go get -x github.com/cogolabs/beyond
CMD ["beyond", "--help"]
