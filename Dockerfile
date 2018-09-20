FROM golang:1.11 as test
WORKDIR /app
COPY go.* /app/
RUN go mod download
COPY *.go /app/
RUN CGO_ENABLED=0 go build -mod=readonly -o /shell-operator ./...

FROM scratch
COPY example/shell-conf.yaml /app/shell-config.yaml
COPY --from=test /shell-operator /shell-operator
ENTRYPOINT ["/shell-operator"]
