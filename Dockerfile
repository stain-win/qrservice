FROM golang:1.14.6-alpine3.12 as builder

COPY go.mod go.sum /go/src/github.com/stain-win/qrservice/
WORKDIR /go/src/github.com/stain-win/qrservice/
RUN go mod download
COPY . /go/src/github.com/stain-win/qrservice/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o dist/qrservice github.com/stain-win/qrservice/cmd/qrservice


FROM alpine

RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/github.com/stain-win/qrservice/dist/qrservice /usr/bin/qrservice

EXPOSE 3200 3200

ENTRYPOINT ["/usr/bin/qrservice"]
