FROM golang:alpine as builder
ENV GOPATH /go
ENV GOFLAGS '-mod=vendor'
ENV GO111MODULE on
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates
RUN adduser -D -g '' tictactoe
WORKDIR ${GOPATH}/src/github.com/codyseavey/tictactoe
COPY . .
RUN chown -R tictactoe assets && chmod -R 777 assets
RUN GOOS=linux GOARCH=386 go build -ldflags="-w -s" -o /go/bin/tictactoe

FROM alpine
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/bin/tictactoe .
COPY --from=builder /go/src/github.com/codyseavey/tictactoe/assets ./assets
USER tictactoe
EXPOSE 8443
ENTRYPOINT ["./tictactoe"]
