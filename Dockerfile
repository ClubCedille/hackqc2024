FROM golang:1.22 AS build

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOPROXY=direct
RUN go build -v -o app .

FROM alpine:3.19
COPY --from=build /go/src/app/app /go/bin/app
COPY --from=build /go/src/app/templates /go/bin/templates
COPY --from=build /go/src/app/docs /go/bin/docs
COPY --from=build /go/src/app/munic_polygons /go/bin/munic_polygons

WORKDIR /go/bin

RUN mkdir -p /go/bin/tmp && chown -R 10001:10001 /go/bin/tmp

ENTRYPOINT ["/go/bin/app"]