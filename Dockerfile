FROM golang:1.11 as test
ENV GO111MODULE auto
WORKDIR /app
COPY go.* /app/
RUN go get
COPY * /app/
RUN go build ./..
