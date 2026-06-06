# Build stage
FROM golang:alpine AS builder

WORKDIR /app

# Instala dependências do sistema necessárias para o build
RUN apk add --no-cache git

# Baixa as dependências do Go
COPY go.mod go.sum ./
RUN go mod download

# Copia o código fonte
COPY . .

# Compila o binário estático
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server/main.go

# Run stage (Imagem super leve)
FROM alpine:latest

WORKDIR /app

# Instala certificados CA para requisições HTTPS (necessário para a Open-Meteo)
RUN apk --no-cache add ca-certificates tzdata

# Copia o binário compilado da etapa anterior
COPY --from=builder /app/main .

# Expor a porta que a aplicação roda
EXPOSE 8080

# Comando para rodar
CMD ["./main"]
