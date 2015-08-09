FROM ubuntu
MAINTAINER Andrew Stuart <andrew.stuart2@gmail.com>

ENTRYPOINT /bible-http-server

EXPOSE 8080
ENV PGPASSWORD=test
ENV BIBLE_PORT=8080

ADD bible-http-server /
