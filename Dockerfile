FROM golang:1.22 AS build
WORKDIR /go/src/app
RUN apk add --no-cache bash
COPY . .
ENV CGO_ENABLED=0 GOOS=linux GOPROXY=direct
RUN go build -v -o app .

FROM scratch
COPY --from=build /go/src/app/app /go/bin/app
COPY --from=build /go/src/app/templates /go/bin/templates
WORKDIR /go/bin
ENTRYPOINT ["/go/bin/app"]