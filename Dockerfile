# Taken from https://github.com/chemidy/smallest-secured-golang-docker-image

FROM golang@sha256:244a736db4a1d2611d257e7403c729663ce2eb08d4628868f9d9ef2735496659 as builder
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

ENV USER=appuser
ENV UID=10001

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}" && \
    mkdir /build && \
    chown -R "${USER}":"${USER}" /build
WORKDIR /build
COPY . .
RUN go get -d -v && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
      -ldflags='-w -s -extldflags "-static"' -a \
      -o shomon .
USER appuser:appuser
ENTRYPOINT [ "./shomon" ]



