FROM golang:alpine
RUN apk add --no-cache git
RUN apk add --no-cache sqlite-libs sqlite-dev
RUN apk add --no-cache build-base
WORKDIR /go/src/app
COPY *.go ./
RUN go-wrapper download
RUN go-wrapper install
RUN apk add alpine-sdk