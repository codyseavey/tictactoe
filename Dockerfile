FROM golang:alpine as builder
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates
RUN adduser -D -g '' tictactoe
WORKDIR $GOPATH/src/github.com/codyseavey/tictactoe
COPY . .
RUN go get -d -v
RUN GOOS=linux GOARCH=386 go build -ldflags="-w -s" -o /go/bin/tictactoe

FROM scratch
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/bin/tictactoe .
USER tictactoe
EXPOSE 8443
ENTRYPOINT ["./tictactoe"]