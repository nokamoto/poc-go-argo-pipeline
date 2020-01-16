FROM golang:1.13.6-alpine3.11 AS build

RUN apk update && apk add git

WORKDIR /src

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY main.go .

RUN go build .

FROM alpine:3.11

COPY --from=build /src/poc-go-argo-pipeline /bin/poc-go-argo-pipeline

ENTRYPOINT [ "/bin/poc-go-argo-pipeline" ]
