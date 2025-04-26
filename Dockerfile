# Этап кэширования модулей
FROM golang:1.24-alpine AS modules
WORKDIR /modules
COPY go.mod go.sum ./
RUN go mod download

# Этап сборки
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Копируем модули из предыдущего этапа
COPY --from=modules /go/pkg /go/pkg

# Копируем исходный код
COPY . .

# Устанавливаем необходимые сертификаты
RUN apk update && apk add --no-cache ca-certificates tzdata
RUN update-ca-certificates

# Устанавливаем переменные сборки
ARG TARGETOS
ARG TARGETARCH
ARG VERSION

# Оптимизированная компиляция с информацией о версии
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -ldflags="-w -s -X main.Version=${VERSION:-dev} -X main.BuildTime=$(date -u +'%Y-%m-%dT%H:%M:%SZ')" \
    -o /bin/app ./cmd/app

# Минимальный финальный образ
FROM scratch

# Добавляем метаданные
LABEL org.opencontainers.image.source="https://github.com/${GITHUB_REPOSITORY}"
LABEL org.opencontainers.image.description="Remnawave Telegram Shop Bot"
LABEL org.opencontainers.image.licenses="MIT"

# Копируем необходимые системные файлы
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Копируем собранное приложение
COPY --from=builder /bin/app /app/app

# Копируем необходимые файлы проекта
COPY --from=builder /app/db /db
COPY --from=builder /app/translations /translations

# Создаем непривилегированного пользователя
USER 1000

# Объявляем порт (документация)
EXPOSE 8080

# Запускаем приложение
CMD ["/app/app"]