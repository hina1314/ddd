FROM golang:1.24.1-alpine AS build

RUN apk add --no-cache ca-certificates tzdata
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
COPY config ./config
COPY db ./db
COPY internal ./internal
COPY token ./token
COPY util ./util

ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w -X main.version=${VERSION}" \
    -o /out/study-api .

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /out/study-api /app/study-api
COPY config/i18n /app/config/i18n

WORKDIR /app
USER 65532:65532
EXPOSE 3000
ENTRYPOINT ["/app/study-api"]
