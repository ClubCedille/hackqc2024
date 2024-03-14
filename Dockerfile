FROM golang:1.22 AS build

WORKDIR /go/src/app
COPY . .
ENV CGO_ENABLED=0 GOOS=linux GOPROXY=direct
RUN go install
RUN go build -v -o app .

FROM alpine:3.19
COPY --from=build /go/src/app/app /go/bin/app
COPY --from=build /go/src/app/templates /go/bin/templates
WORKDIR /go/bin
RUN mkdir -p /go/bin/tmp

ENTRYPOINT ["/go/bin/app"]