FROM golang:latest as build_env

RUN mkdir /app 
ADD . /app
COPY . /app
WORKDIR /app 

ENV GOOS=linux GOARCH=amd64 CGO_ENABLED=0
RUN go build -o /app/server ./cmd/serverd

FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM alpine:latest
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build_env /app/server /bin/server

ENTRYPOINT ["server"]
