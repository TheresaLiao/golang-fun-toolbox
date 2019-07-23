FROM golang:1.11.2-alpine

RUN apk update && apk upgrade
RUN apk add --no-cache git

RUN go get -u github.com/gin-gonic/gin
RUN go get -u github.com/op/go-logging


RUN mkdir -p /app
WORKDIR /app

ADD . /app

COPY ca.crt /var/run/secrets/kubernetes.io/serviceaccount/
RUN go build -o main
EXPOSE 80
CMD /app/main
