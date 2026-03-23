FROM golang:1.26

# Install air
RUN go install github.com/air-verse/air@latest

ENV GO_ENV=development

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 5000

CMD ["air", "-c", ".air.toml"]
