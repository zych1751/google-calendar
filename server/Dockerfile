FROM golang:alpine AS builder

RUN apk update
RUN apk add -U --no-cache ca-certificates && update-ca-certificates

WORKDIR /usr/src/app
COPY . .

RUN go get -u -v ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-w -s' -o main .


FROM scratch

EXPOSE 9000
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/src/app .

CMD ["./main"]