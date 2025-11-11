FROM golang:1.25.4

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN apt-get update && \
    apt-get install -y --no-install-recommends postgresql-client curl bash && \
    rm -rf /var/lib/apt/lists/* && \
    curl -fsSL https://github.com/pressly/goose/releases/download/v3.16.0/goose_linux_x86_64 \
        -o /usr/local/bin/goose && \
    chmod +x /usr/local/bin/goose

EXPOSE 8080

CMD ["go", "run", "./cmd/main.go"]