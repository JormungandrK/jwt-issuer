### Multi-stage build
FROM golang:1.13.5-alpine3.10 as build

RUN apk --no-cache add git curl openssh

COPY . /go/src/github.com/Microkubes/jwt-issuer

RUN cd /go/src/github.com/Microkubes/jwt-issuer && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install


### Main
FROM alpine:3.10

COPY --from=build /go/src/github.com/Microkubes/jwt-issuer/config.json /config.json
COPY --from=build /go/bin/jwt-issuer /jwt-issuer

EXPOSE 8080

CMD ["/jwt-issuer"]
