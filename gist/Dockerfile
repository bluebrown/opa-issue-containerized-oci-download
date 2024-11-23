FROM docker.io/library/golang
WORKDIR /go/src
COPY go.mod go.sum ./
RUN go mod download -x
COPY *.go ./
RUN go build -v -o app .

FROM docker.io/library/debian
RUN apt-get update && apt-get install -y ca-certificates
COPY --from=0 /go/src/app /opt/acme/app
ENTRYPOINT ["/opt/acme/app"]

