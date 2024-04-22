FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go mod tidy

COPY . .

RUN go build -o scope_test ./cmd

ENV PORT=8080

CMD ["./scope_test", "-port=$PORT"]