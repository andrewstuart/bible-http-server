FROM golang:latest 
MAINTAINER Andrew Stuart <andrew.stuart2@gmail.com>

EXPOSE 8080
ENV PGPASSWORD=test
ENV BIBLE_PORT=8080
ENV VIRTUAL_HOST=bible-api,bible-api.astuart.co
CMD ["/go/bin/bible-http-server"]

RUN mkdir -p /go/src/github.com/andrewstuart/bible-http-server
ADD . /go/src/github.com/andrewstuart/bible-http-server
RUN go get github.com/andrewstuart/bible-http-server
