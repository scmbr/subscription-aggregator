FROM golang:1.24
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

EXPOSE 8000

CMD ["sh",  "go run ./cmd/app"]
