FROM golang:alpine AS builder

ENV GOTOOLCHAIN=auto

WORKDIR /app

# Copia os módulos compartilhados
COPY pkg ./pkg

# Copia e constrói o serviço
COPY backend-catalog ./backend-catalog
WORKDIR /app/backend-catalog
RUN go mod download
RUN go build -o /app/bin/api ./cmd/api

FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/bin/api .

EXPOSE 8086

CMD ["./api"]
