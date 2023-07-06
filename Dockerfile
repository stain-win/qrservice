FROM golang:1.20.0-alpine as builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -mod=readonly -v -o qrgen ./cmd/qrservice


FROM alpine:3

RUN apk update \
      && apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /app/qrgen /qrgen

ENTRYPOINT ["/qrgen"]
