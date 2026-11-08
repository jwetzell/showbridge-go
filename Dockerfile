ARG GO_VERSION=1.27.3

FROM node:26-alpine AS webui-builder
WORKDIR /webui
COPY webui .
RUN npm ci && npm run build -- --base-href=/ui/

FROM golang:${GO_VERSION}-alpine AS build
RUN apk --no-cache add ca-certificates tzdata git build-base alsa-lib-dev
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=webui-builder /webui/dist/webui ./api/webui
RUN CGO_ENABLED=1 GOOS=linux go build -tags embed_webui ./cmd/showbridge

FROM alpine:3.24.4
RUN apk --no-cache add ca-certificates tzdata alsa-lib libgcc libstdc++
WORKDIR /app
COPY --from=build /build/showbridge /app/showbridge
COPY config.min.yaml /data/config.yaml
ENV SHOWBRIDGE_CONFIG=/data/config.yaml
ENTRYPOINT [ "/app/showbridge" ]