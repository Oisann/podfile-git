FROM golang:1.10-alpine as builder

WORKDIR /go/src/github.com/oisann/
COPY main.go main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM scratch
RUN apk --update add --no-cache bash gawk sed grep bc coreutils git openssh && \
    rm -rf /var/lib/apt/lists/* && \
    rm /var/cache/apk/*
WORKDIR /git
COPY --from=builder /go/src/github.com/oisann/main .
ENTRYPOINT ["./main"]