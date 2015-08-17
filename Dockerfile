FROM golang:latest 
MAINTAINER Andrew Stuart <andrew.stuart2@gmail.com>

EXPOSE 8080
ENV PGPASSWORD=test
ENV BIBLE_PORT=8080
ENV VIRTUAL_HOST=bible-api,bible-api.astuart.co
CMD ["/$GOPATH/bin/bible-http-server"]

ENV APP_GETPATH=github.com/andrewstuart/bible-http-server
ENV APP_PATH=/$GOPATH/src/$APP_GETPATH
ADD . $APP_PATH
RUN go get $APP_GETPATH
