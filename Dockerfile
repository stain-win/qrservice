FROM golang:1.14.6-alpine3.12 as builder

COPY go.mod go.sum /go/src/gitlab.com/stain-win/qrservice/
WORKDIR /go/src/gitlab.com/stain-win/qrservice/
RUN go mod download
COPY . /go/src/gitlab.com/stain-win/qrservice/
RUN ./build.sh


FROM alpine

RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/gitlab.com/stain-win/qrservice/dist/qrservice /usr/bin/qrservice

EXPOSE 8080 8080

ENTRYPOINT ["/usr/bin/qrservice"]
