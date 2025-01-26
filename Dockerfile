# Use a imagem oficial do Golang como builder
FROM golang:1.23.4 AS builder

WORKDIR /app

# Copie os módulos Go e instale as dependências
COPY go.mod go.sum ./
RUN go mod download

# Copie o código fonte
COPY . .

# Compile a aplicação
RUN go build -o rate-limiter ./cmd/main.go

# Use uma imagem mínima para rodar a aplicação
FROM alpine:latest

WORKDIR /root/

# Copie a aplicação compilada do builder
COPY --from=builder /app/rate-limiter .

# Copie o arquivo .env
COPY .env .

# Comando para rodar a aplicação
CMD ["./rate-limiter"]
