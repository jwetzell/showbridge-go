ARG GO_VERSION=1.25.5
FROM golang:${GO_VERSION}-alpine AS build
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/showbridge

FROM scratch
WORKDIR /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /build/showbridge /app/showbridge
ENTRYPOINT [ "/app/showbridge" ]