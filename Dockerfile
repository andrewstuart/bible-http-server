FROM golang:latest 
MAINTAINER Andrew Stuart <andrew.stuart2@gmail.com>

EXPOSE 8080
ENV PGPASSWORD=test
ENV BIBLE_PORT=8080

RUN mkdir -p /go/src/github.com/andrewstuart/bible-http-server
ADD . /go/src/github.com/andrewstuart/bible-http-server
RUN go get github.com/andrewstuart/bible-http-server
CMD ["/go/bin/bible-http-server"]
