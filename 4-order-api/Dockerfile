# ---- build stage ----
FROM golang:1.25-alpine AS build

WORKDIR /src

# Кэшируем зависимости отдельным слоем
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники и собираем статический бинарник
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./Cmd

# ---- runtime stage ----
FROM alpine:3.20

# ca-certificates на случай исходящих TLS-запросов (SMS-провайдер и т.п.)
RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=build /out/server /app/server

EXPOSE 8081

ENTRYPOINT ["/app/server"]
