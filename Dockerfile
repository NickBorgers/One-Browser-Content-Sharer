FROM golang:latest AS builder
ADD app.go /go/src/project/app.go
WORKDIR /go/src/project/
RUN go get ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /executable /go/src/project/app.go

FROM scratch
COPY --from=builder /executable .
COPY ./static /static
ENTRYPOINT ["./executable"]
