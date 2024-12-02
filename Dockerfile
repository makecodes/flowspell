# Etapa de construção
FROM golang:1.22.9 AS builder

# Define o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copia os arquivos go.mod e go.sum para o diretório de trabalho
COPY go.mod go.sum ./

# Baixa as dependências do módulo
RUN go mod download

# Copia o restante dos arquivos do projeto para o diretório de trabalho
COPY . .

# Compila o aplicativo
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o flowspell main.go

# Etapa final
FROM gcr.io/distroless/base-debian10

# Define o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copia o binário compilado da etapa de construção para a etapa final
COPY --from=builder /app/flowspell .

# Define o comando de inicialização do contêiner
CMD ["./flowspell"]
