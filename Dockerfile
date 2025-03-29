# FROM golang:1.16-buster AS build

FROM golang:1.16-alpine

ARG MODULE_PATH=/go/src/bitbucket.org/ssinbeti/commission-module

RUN mkdir -p ${MODULE_PATH}
WORKDIR ${MODULE_PATH}

ADD . .

RUN go mod tidy
RUN go mod vendor

RUN go build -o ./main ./cmd/rest-server/main.go

CMD ["./main"]