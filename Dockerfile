# FROM golang:1.16-buster AS build

FROM golang:1.19.0-alpine3.15

ARG MODULE_PATH=/go/src/bitbucket.org/ssinbeti/commissionn-module

RUN mkdir -p ${MODULE_PATH}
WORKDIR ${MODULE_PATH}

ADD . .

RUN go mod tidy
RUN go mod vendor

RUN go build -o ./main ./cmd/rest-server/main.go

CMD ["./main"]