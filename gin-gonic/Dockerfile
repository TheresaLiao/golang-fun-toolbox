FROM golang:1.11.2-alpine

RUN apk update && apk upgrade
RUN apk add --no-cache git

RUN go get -u github.com/gin-gonic/gin

RUN mkdir -p /app
WORKDIR /app

ADD . /app

RUN go build -o main
EXPOSE 80
CMD /app/main
