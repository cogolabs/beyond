FROM golang
ADD . /go/src/github.com/cogolabs/transcend
RUN go get -x github.com/cogolabs/transcend
CMD ["transcend", "--help"]
