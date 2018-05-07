FROM golang:1.10-alpine as builder

WORKDIR /go/src/github.com/oisann/
COPY main.go main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine/git:latest
WORKDIR /git
COPY --from=builder /go/src/github.com/oisann/main /bin/
ENTRYPOINT ["main"]
