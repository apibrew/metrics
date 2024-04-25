FROM golang:1.21-alpine as buildenv

WORKDIR /app/
#RUN apk add build-base

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

FROM buildenv as lint
WORKDIR /app/

RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2

COPY pkg pkg

RUN golangci-lint run --timeout 5m --verbose

FROM buildenv as test
WORKDIR /app/

# setup database
RUN apk add postgresql
RUN mkdir /run/postgresql
RUN chown postgres:postgres /run/postgresql/

COPY pkg pkg

RUN sh /app/pkg/test/run-tests.sh
