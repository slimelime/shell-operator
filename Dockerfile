FROM golang:1.11.0 as base
RUN go get -u github.com/golang/dep/cmd/dep
WORKDIR /go/src/github.com/MYOB-Technology/shell-operator

FROM golang:1.11.0 as build
WORKDIR /go/src/github.com/MYOB-Technology/shell-operator
COPY --from=base /go/bin/dep /go/bin/dep
COPY Gopkg.* /go/src/github.com/MYOB-Technology/shell-operator/
RUN dep ensure -v -vendor-only
COPY . /go/src/github.com/MYOB-Technology/shell-operator/
RUN CGO_ENABLED=0 go build -o /shell-operator main.go

FROM scratch
COPY example/shell-conf.yaml /app/shell-config.yaml
COPY --from=build /shell-operator /shell-operator
ENTRYPOINT ["/shell-operator"]
