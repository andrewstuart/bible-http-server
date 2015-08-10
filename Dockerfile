FROM golang:latest 
MAINTAINER Andrew Stuart <andrew.stuart2@gmail.com>

EXPOSE 8080
ENV PGPASSWORD=test
ENV BIBLE_PORT=8080

RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
RUN go build -o main . 
CMD ["/app/main"]


