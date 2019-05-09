FROM golang:1.12 AS builder

EXPOSE 8080

WORKDIR /go/src/shorty
COPY . .

RUN go get -d -v ./...
RUN go build -o /shorty -ldflags "-linkmode external -extldflags -static" -a *.go

FROM scratch
COPY --from=builder /shorty /shorty
CMD ["/shorty"]