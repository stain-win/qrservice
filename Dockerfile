FROM golang:1.20.0-alpine as builder

ARG TARGETOS
ARG TARGETARCH

COPY go.mod go.sum ./qrservice/
WORKDIR /app
RUN go mod download
COPY . /app
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -mod=readonly -a -installsuffix cgo -o qrgen


FROM alpine:3

RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /app/qrgen /usr/bin/qrservice

EXPOSE 3200 3200

ENTRYPOINT ["/usr/bin/qrservice"]
