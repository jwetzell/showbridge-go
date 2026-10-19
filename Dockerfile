ARG GO_VERSION=1.27.2

FROM node:24-alpine AS webui-builder
WORKDIR /webui
COPY webui .
RUN npm ci && npm run build -- --base-href=/ui/

FROM golang:${GO_VERSION}-alpine AS build
RUN apk --no-cache add ca-certificates tzdata git
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=webui-builder /webui/dist/webui ./internal/api/webui
RUN CGO_ENABLED=0 GOOS=linux go build -tags embed_webui ./cmd/showbridge

FROM scratch
WORKDIR /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /build/showbridge /app/showbridge
COPY config.min.yaml /data/config.yaml
ENV SHOWBRIDGE_CONFIG=/data/config.yaml
ENTRYPOINT [ "/app/showbridge" ]