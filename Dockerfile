FROM golang:1.9

RUN go get github.com/ToQoz/gopwt/...
RUN go get github.com/gopherjs/gopherjs

EXPOSE 5000
WORKDIR "/go/src/github.com/ToQoz/gopwt.toqoz.net"

ADD web.go  web.go
ADD sandbox sandbox
ADD statics statics

RUN go get ./...

CMD ["go", "run", "web.go"]
