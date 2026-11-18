# stage 1: build the go application
FROM golang:1.27.3-alpine3.24 AS builder

WORKDIR /app

RUN apk add --no-cache make

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN make docker-build

# satge 2: build the apline application

FROM alpine:3.23

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/bin/main ./app

EXPOSE 8080

CMD [ "./main" ]
