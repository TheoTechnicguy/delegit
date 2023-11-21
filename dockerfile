FROM golang:1.21 AS build

WORKDIR /opt/delegit

COPY go.mod go.sum ./
RUN go mod download

COPY *.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o delegit

FROM gcr.io/distroless/base-debian11 AS main

WORKDIR /opt/delegit
COPY --from=build /opt/delegit/delegit .

# ARG BASE_DIGEST

LABEL org.opencontainers.image.authors "Delegit Community"
# LABEL org.opencontainers.image.url "https://git.licolas.net/delegit/delegit"
LABEL org.opencontainers.image.documentation "https://git.licolas.net/delegit/delegit"
# LABEL org.opencontainers.image.source "https://git.licolas.net/delegit/delegit"
LABEL org.opencontainers.image.version ${VERSION}
LABEL org.opencontainers.image.vendor "Delegit Community"
LABEL org.opencontainers.image.licenses "Apache-2.0"
# LABEL org.opencontainers.image.ref.name
LABEL org.opencontainers.image.title "Delegit Server Application"
LABEL org.opencontainers.image.description "This image contains the proof-of-concept Delegit HTTP Server"
# LABEL org.opencontainers.image.base.digest ${BASE_DIGEST}
LABEL org.opencontainers.image.base.name "gcr.io/distroless/base-debian11:latest-arm64"

EXPOSE 41990

ENTRYPOINT ["/opt/delegit/delegit"]
