FROM golang:alpine as build
LABEL authors="snappey"

RUN apk add -U --no-cache ca-certificates

WORKDIR /app

COPY . .

RUN go get

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./brb

FROM scratch

COPY --from=build /app/brb /brb
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/brb"]