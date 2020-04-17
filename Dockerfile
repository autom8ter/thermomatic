FROM golang:1.12.9-alpine3.10 as build-env

RUN apk add git
RUN mkdir /thermomatic
RUN apk --update add ca-certificates
WORKDIR /thermomatic
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/thermomatic
FROM scratch
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build-env /go/bin/thermomatic /go/bin/thermomatic
WORKDIR /thermomatic
ENTRYPOINT ["/go/bin/thermomatic"]