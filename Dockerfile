# Use a imagem oficial do Golang como imagem de construção.
FROM golang:1.22 AS builder

# Define o diretório de trabalho dentro do contêiner.
WORKDIR /app

# Copia o módulo go e o arquivo de soma para o diretório de trabalho.
COPY go.mod .
COPY go.sum .

# Baixa as dependências de Go.
RUN go mod download

# Copia todos os arquivos do código fonte para o diretório de trabalho.
COPY . .

# Altera o diretório de trabalho para o diretório onde o main.go está localizado.
WORKDIR /app/cmd

# Compila a aplicação Go para um binário estático.
# Certifique-se de ajustar o caminho do main.go se necessário.
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp .

# Use a imagem do alpine para a imagem final.
FROM alpine:latest  
WORKDIR /root/

# Instala certificados CA para permitir que a aplicação faça chamadas HTTPS.
RUN apk --no-cache add ca-certificates

# Copia o binário estático para a imagem final.
COPY --from=builder /app/cmd/myapp .

# Executa o binário.
CMD ["./myapp"]


