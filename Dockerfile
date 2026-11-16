ARG GO_VERSION=1.27.3


FROM golang:${GO_VERSION}-alpine AS build
RUN apk --no-cache add ca-certificates tzdata git build-base alsa-lib-dev
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build ./cmd/showbridge

FROM alpine:3.24.4
RUN apk --no-cache add ca-certificates tzdata alsa-lib libgcc libstdc++
WORKDIR /app
COPY --from=build /build/showbridge /app/showbridge
COPY config.min.yaml /data/config.yaml
ENV SHOWBRIDGE_CONFIG=/data/config.yaml
ENTRYPOINT [ "/app/showbridge" ]