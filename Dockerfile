FROM golang:1.22.0-bullseye

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o server ./cmd/main.go


CMD ["./server"]
